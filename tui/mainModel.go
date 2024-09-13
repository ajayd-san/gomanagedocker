package tui

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	teadialog "github.com/ajayd-san/teaDialog"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"

	"github.com/ajayd-san/gomanagedocker/service"
	"github.com/ajayd-san/gomanagedocker/tui/components/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.design/x/clipboard"

	"github.com/ajayd-san/gomanagedocker/service/dockercmd"
	it "github.com/ajayd-san/gomanagedocker/service/types"
	"github.com/ajayd-san/gomanagedocker/tui/components"
)

// dimension ratios for infobox
const infoBoxWidthRatio = 0.55
const infoBoxHeightRatio = 0.6

// type for denoting the tab order
type tabId int

type TickMsg time.Time

// when the main loop gets a variable of this type, it preloads all required tabs.
// TODO: Should change type from `int` to `struct{}`
type preloadObjects int

// instruction for preload contaienr sizes, when program starts.
type preloadSizeMap struct{}

type ContainerSize struct {
	sizeRw int64
	rootFs int64
}

// this type holds info about container sizes and is meant to be used concurrently.
type ContainerSizeManager struct {
	sizeMap map[string]ContainerSize
	mu      *sync.Mutex
}

type MainModel struct {
	serviceKind         it.ServiceType
	dockerClient        service.Service
	Tabs                []string
	TabContent          []listModel
	activeTab           tabId
	width               int
	height              int
	windowTooSmall      bool
	displayInfoBox      bool
	windowtoosmallModel WindowTooSmallModel
	keymap              KeyMap
	//  handles navigation keymap generation
	navKeymap help.Model
	// handles tab specific keymap generation, i have no idea why I named it `helpGen`
	helpGen      help.Model
	showDialog   bool
	activeDialog tea.Model
	// we use this to cancel dialog ops when we exit from them
	dialogOpCancel context.CancelFunc
	// this maintains a map of container image sizes
	containerSizeTracker ContainerSizeManager
	// this maps imageIds to imageNames (for legibility)
	imageIdToNameMap map[string]string
	// we use this error channel to report error for possibly long running tasks, like pruning
	possibleLongRunningOpErrorChan chan error

	// Channels for sending and receiving notifications, we use these to update list status messages
	notificationChan chan notificationMetadata

	// only used to store fatal error before quitting with non-zero exit code
	exitError error
}

// this ticker enables us to update Docker lists items every 500ms (unless set to different value in config)
func doUpdateObjectsTick() tea.Cmd {
	return tea.Tick(CONFIG_POLLING_TIME, func(t time.Time) tea.Msg { return TickMsg(t) })
}

func (m MainModel) Init() tea.Cmd {
	// check if Docker is alive, if not, exit
	err := m.dockerClient.Ping()
	if err != nil {
		earlyExitErr = fmt.Errorf("Error connecting to Docker daemon.\nInfo: %w\n", err)
		return tea.Quit
	}

	// initialize clipboard
	// TODO: handle error
	clipboard.Init()
	// this command enables loading tab contents a head of time, so there is no load time while switching tabs
	preloadCmd := func() tea.Msg { return preloadObjects(0) }
	preloadSize := func() tea.Msg { return preloadSizeMap{} }
	return tea.Batch(preloadCmd, preloadSize, doUpdateObjectsTick())
}

// Initializes and returns a new Model instance.
func NewModel(client service.Service, serviceType it.ServiceType) MainModel {
	contents := make([]listModel, len(CONFIG_TAB_ORDERING))
	keymap := NewKeyMap(serviceType)

	for tabid, tabName := range CONFIG_TAB_ORDERING {
		var objectHelp help.KeyMap
		var objectHelpBulk help.KeyMap

		switch tabName {
		case "images":
			objectHelp = keymap.image
			objectHelpBulk = keymap.imageBulk
		case "containers":
			objectHelp = keymap.container
			objectHelpBulk = keymap.containerBulk
		case "volumes":
			objectHelp = keymap.volume
			objectHelpBulk = keymap.volumeBulk
		case "pods":
			objectHelp = keymap.pods
			objectHelpBulk = keymap.podsBulk
		}
		contents[tabid] = InitList(tabId(tabid), objectHelp, objectHelpBulk)
	}

	firstTab := contents[0].tabKind

	helper := help.New()
	helper.FullSeparator = " • "
	helper.ShowAll = true

	NavKeymap := help.New()
	NavKeymap.FullSeparator = " • "
	NavKeymap.ShowAll = true

	return MainModel{
		serviceKind:                    serviceType,
		dockerClient:                   client,
		Tabs:                           CONFIG_TAB_ORDERING,
		TabContent:                     contents,
		displayInfoBox:                 true,
		windowtoosmallModel:            MakeNewWindowTooSmallModel(),
		possibleLongRunningOpErrorChan: make(chan error, 10),
		keymap:                         keymap,
		helpGen:                        helper,
		navKeymap:                      NavKeymap,
		activeTab:                      firstTab,
		containerSizeTracker: ContainerSizeManager{
			sizeMap: make(map[string]ContainerSize),
			mu:      &sync.Mutex{},
		},
		imageIdToNameMap: make(map[string]string),
		notificationChan: make(chan notificationMetadata, 10),
	}
}

// goMangeDocker main update loop
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd

notificationLoop:
	for {
		select {
		case notifcation := <-m.notificationChan:
			cmd := (m.TabContent[notifcation.listId].list).NewStatusMessage(notifcation.msg)
			cmds = append(cmds, cmd)
		default:
			break notificationLoop
		}
	}

	//check if error exists on error channel, if yes show the error in new dialog
	select {
	case newErr := <-m.possibleLongRunningOpErrorChan:
		if newErr != nil {
			m.showDialog = true
			m.activeDialog = teadialog.NewErrorDialog(newErr.Error(), m.width)
		}
	default:
	}

	// INFO: if m.showDialog is true, then hijack all keyinputs and forward them to the dialog
	if m.showDialog {

		update, cmd := m.activeDialog.Update(msg)
		m.activeDialog = update

		// if keymsg is <Esc> then close dialog
		if msg, ok := msg.(tea.KeyMsg); ok && key.Matches(msg, m.keymap.navigation.Back) {
			if m.dialogOpCancel != nil {
				m.dialogOpCancel()
				// this might be required, in the future
				// m.dialogOpCancel = nil
			}
			m.showDialog = false
		}

		cmds = append(cmds, cmd)
	}

	switch assertedMsg := msg.(type) {
	// fetches container size info in a separate go routine
	case preloadSizeMap:
		go m.prepopulateContainerSizeMapConcurrently()

	// preloads all tabs, so no delay in displaying objects when first changing tabs
	case preloadObjects:
		//FIXME: use a range on TAB_ORDERING to preload lists

		for tabid := range CONFIG_TAB_ORDERING {
			m = m.updateContent(tabId(tabid))
		}

	case TickMsg:
		m = m.updateContent(m.activeTab)

		cmds = append(cmds, doUpdateObjectsTick())

	case tea.WindowSizeMsg:
		// show windowtoosmallModel if window dimentions are too small
		if assertedMsg.Height < 25 || assertedMsg.Width < 65 {
			m.windowTooSmall = true
			temp, _ := m.windowtoosmallModel.Update(assertedMsg)
			m.windowtoosmallModel = temp.(WindowTooSmallModel)
		} else {
			m.windowTooSmall = false
		}

		// toggle info box if window size goes under a certain threshold
		if assertedMsg.Height <= 31 || assertedMsg.Width < 105 {
			m.displayInfoBox = false
			listWidthRatio = listWidthRatioWithOutInfoBox
		} else {
			listWidthRatio = listWidthRatioWithInfoBox
			m.displayInfoBox = true
		}

		m.width = assertedMsg.Width
		m.height = assertedMsg.Height

		KeymapAvailableWidth = assertedMsg.Width - 10

		windowStyle = windowStyle.
			Width(m.width - listDocStyle.GetHorizontalFrameSize() - 2).
			Height(m.height - listDocStyle.GetVerticalFrameSize() - 3)

		dialogContainerStyle = dialogContainerStyle.Width(assertedMsg.Width).Height(assertedMsg.Height)

		// dynamically resizes the dimentions of infobox depending on the window size
		moreInfoStyle = moreInfoStyle.Width(int(infoBoxWidthRatio * float64(m.width)))
		moreInfoStyle = moreInfoStyle.Height(int(infoBoxHeightRatio * float64(m.height)))

		m.helpGen.Width = assertedMsg.Width - 20
		m.navKeymap.Width = assertedMsg.Width - 20

		// change list dimensions when window size changes
		// TODO: change width
		for index := range m.TabContent {
			listM, _ := m.TabContent[index].Update(assertedMsg)
			m.TabContent[index] = listM.(listModel)
		}

	case tea.KeyMsg:
		if !m.getActiveList().SettingFilter() && !m.showDialog {
			switch {
			case key.Matches(assertedMsg, m.keymap.navigation.Quit):
				return m, tea.Quit
			case key.Matches(assertedMsg, m.keymap.navigation.NextTab):
				m.nextTab()
			case key.Matches(assertedMsg, m.keymap.navigation.PrevTab):
				m.prevTab()
			case key.Matches(assertedMsg, m.keymap.navigation.Select):
				msg = itemSelect{}
				break
			case key.Matches(assertedMsg, m.keymap.navigation.Back):
				if m.getActiveTab().inBulkMode() {
					msg = clearSelection{}
					break
				}
			}

			if m.activeTab == IMAGES {
				switch {
				case key.Matches(assertedMsg, m.keymap.image.Run):
					curItem := m.getSelectedItem()

					if curItem != nil && !m.isCurrentTabInBulkMode() {
						imageInfo := curItem.(imageItem)
						storage := map[string]string{"ID": imageInfo.GetId()}
						m.activeDialog = getRunImageDialog(storage)
						m.showDialog = true
						cmds = append(cmds, m.activeDialog.Init())
						// op := runImage(m.dockerClient, imageInfo, m.activeTab, m.notificationChan)
						// go m.runBackground(op)
					}

				case key.Matches(assertedMsg, m.keymap.image.Delete):
					curItem := m.getSelectedItem()

					if curItem != nil && !m.isCurrentTabInBulkMode() {
						imageId := curItem.(dockerRes).GetId()
						storage := map[string]string{"ID": imageId}
						m.activeDialog = getRemoveImageDialog(storage)
						m.showDialog = true

						cmds = append(cmds, m.activeDialog.Init())
					}

				case key.Matches(assertedMsg, m.keymap.image.DeleteForce):
					items := m.getSelectedItems()

					deleteOpts := it.RemoveImageOptions{
						Force:   true,
						NoPrune: true,
					}
					op := imageDeleteBulk(
						m.dockerClient,
						items,
						deleteOpts,
						m.activeTab,
						m.notificationChan,
						m.possibleLongRunningOpErrorChan,
					)
					go m.runBackground(op)

					cmds = append(cmds, clearSelectionCmd())

				case key.Matches(assertedMsg, m.keymap.image.Prune):
					if !m.isCurrentTabInBulkMode() {
						m.activeDialog = getPruneImagesDialog(make(map[string]string))
						m.showDialog = true
						cmds = append(cmds, m.activeDialog.Init())
					}

				case key.Matches(assertedMsg, m.keymap.image.Scout):
					curItem := m.getSelectedItem()
					if curItem != nil && !m.isCurrentTabInBulkMode() {
						dockerRes := curItem.(dockerRes)
						imageInfo := dockerRes.(imageItem)
						imageName := imageInfo.RepoTags[0]

						ctx, cancel := context.WithCancel(context.Background())
						m.dialogOpCancel = cancel

						f := func() (*dockercmd.ScoutData, error) {

							dockerClient := m.dockerClient.(*dockercmd.DockerClient)
							scoutData, err := dockerClient.ScoutImage(ctx, imageName)

							if err != nil {
								m.possibleLongRunningOpErrorChan <- err
							}

							return scoutData, err
						}

						m.activeDialog = getImageScoutDialog(f)
						m.showDialog = true
						cmds = append(cmds, m.activeDialog.Init())
					}
				case key.Matches(assertedMsg, m.keymap.image.CopyId):
					currentItem := m.getSelectedItem()

					if currentItem != nil && !m.isCurrentTabInBulkMode() {
						dres := currentItem.(dockerRes)
						op := copyIdToClipboard(dres, m.activeTab, m.notificationChan)
						op()
					}

				case key.Matches(assertedMsg, m.keymap.image.RunAndExec):
					currentItem := m.getSelectedItem()

					if currentItem != nil && !m.isCurrentTabInBulkMode() {
						dres := currentItem.(dockerRes)
						id := dres.GetId()

						// config := container.Config{
						// 	AttachStdin:  true,
						// 	AttachStdout: true,
						// 	AttachStderr: true,
						// 	Tty:          true,
						// 	Image:        id,
						// }

						config := it.ContainerCreateConfig{
							ImageId: id,
						}

						containerId, err := m.dockerClient.RunImage(config)
						if err != nil {
							m.activeDialog = teadialog.NewErrorDialog(err.Error(), m.width)
							m.showDialog = true
						}

						cmd := m.dockerClient.ExecCmd(*containerId)
						cmds = append(cmds, tea.ExecProcess(cmd, func(err error) tea.Msg {
							m.possibleLongRunningOpErrorChan <- err
							return nil
						}))
					}

				case key.Matches(assertedMsg, m.keymap.image.Build):
					if !m.isCurrentTabInBulkMode() {
						m.activeDialog = getBuildImageDialog(make(map[string]string))
						m.showDialog = true
						cmds = append(cmds, m.activeDialog.Init())
					}
				}

			} else if m.activeTab == CONTAINERS {
				switch {
				case key.Matches(assertedMsg, m.keymap.container.ToggleListAll):
					toggleListAllContainers(m.dockerClient, m.activeTab, m.notificationChan)

				case key.Matches(assertedMsg, m.keymap.container.ToggleStartStop):
					selectedItems := m.getSelectedItems()

					op := toggleStartStopContainer(m.dockerClient, selectedItems, m.activeTab, m.notificationChan, m.possibleLongRunningOpErrorChan)
					go m.runBackground(op)
					cmds = append(cmds, clearSelectionCmd())

				case key.Matches(assertedMsg, m.keymap.container.TogglePause):
					selectedItems := m.getSelectedItems()

					op := togglePauseResumeContainer(m.dockerClient, selectedItems, m.activeTab, m.notificationChan, m.possibleLongRunningOpErrorChan)
					go m.runBackground(op)
					cmds = append(cmds, clearSelectionCmd())

				case key.Matches(assertedMsg, m.keymap.container.Restart):
					selectedItems := m.getSelectedItems()

					op := toggleRestartContainer(m.dockerClient, selectedItems, m.activeTab, m.notificationChan, m.possibleLongRunningOpErrorChan)
					go m.runBackground(op)
					cmds = append(cmds, clearSelectionCmd())

				case key.Matches(assertedMsg, m.keymap.container.Delete):
					curItem := m.getSelectedItem()
					if curItem != nil && !m.isCurrentTabInBulkMode() {
						containerInfo := curItem.(dockerRes)
						dialog := getRemoveContainerDialog(map[string]string{"ID": containerInfo.GetId()})
						m.activeDialog = dialog
						m.showDialog = true
						cmds = append(cmds, m.activeDialog.Init())
					}

				case key.Matches(assertedMsg, m.keymap.container.DeleteForce):
					selectedItems := m.getSelectedItems()

					deleteOpts := it.ContainerRemoveOpts{
						RemoveVolumes: false,
						RemoveLinks:   false,
						Force:         true,
					}

					op := containerDeleteBulk(
						m.dockerClient,
						selectedItems,
						deleteOpts,
						m.activeTab,
						m.notificationChan,
						m.possibleLongRunningOpErrorChan,
					)
					go m.runBackground(op)

					cmds = append(cmds, clearSelectionCmd())

				case key.Matches(assertedMsg, m.keymap.container.Prune):
					if !m.isCurrentTabInBulkMode() {
						m.activeDialog = getPruneContainersDialog(make(map[string]string))
						m.showDialog = true
						cmds = append(cmds, m.activeDialog.Init())
					}

				case key.Matches(assertedMsg, m.keymap.container.Exec):
					curItem := m.getSelectedItem()
					if curItem != nil && !m.isCurrentTabInBulkMode() {
						container := curItem.(containerItem)

						if container.getState() != "running" {
							//get the article for correct grammar
							var article string
							if container.getState() == "exited" {
								article = "an"
							} else {
								article = "a"
							}

							m.activeDialog = teadialog.NewErrorDialog(
								fmt.Sprintf("Cannot exec into %s %s container", article, container.getState()),
								m.width,
							)
							m.showDialog = true
							cmds = append(cmds, m.activeDialog.Init())
						} else {
							containerId := container.GetId()
							// execs into the default shell of the container (got from lazydocker)
							cmd := m.dockerClient.ExecCmd(containerId)

							cmds = append(cmds, tea.ExecProcess(cmd, func(err error) tea.Msg {
								m.possibleLongRunningOpErrorChan <- err
								return nil
							}))
						}
					}

				case key.Matches(assertedMsg, m.keymap.container.CopyId):
					currentItem := m.getSelectedItem()

					if currentItem != nil && !m.isCurrentTabInBulkMode() {
						object := currentItem.(dockerRes)
						op := copyIdToClipboard(object, m.activeTab, m.notificationChan)

						op()
					}

				case key.Matches(assertedMsg, m.keymap.container.ShowLogs):
					currentItem := m.getSelectedItem()

					if currentItem != nil && !m.isCurrentTabInBulkMode() {
						dres := currentItem.(containerItem)
						if dres.State == "running" {
							id := dres.GetId()
							cmd := m.dockerClient.LogsCmd(id)
							cmds = append(cmds, tea.ExecProcess(cmd, func(err error) tea.Msg {
								if err.Error() != "exit status 1" {
									m.possibleLongRunningOpErrorChan <- err
								}
								return nil
							}))
						}
					}
				}

			} else if m.activeTab == VOLUMES {
				switch {
				case key.Matches(assertedMsg, m.keymap.volume.Prune):
					curItem := m.getSelectedItem()
					if curItem != nil && !m.isCurrentTabInBulkMode() {
						volumeId := curItem.(dockerRes).GetId()
						m.activeDialog = getPruneVolumesDialog(map[string]string{"ID": volumeId})
						m.showDialog = true
						cmds = append(cmds, m.activeDialog.Init())
					}

				case key.Matches(assertedMsg, m.keymap.volume.Delete):

					curItem := m.getSelectedItem()

					if curItem != nil && !m.isCurrentTabInBulkMode() {
						volumeId := curItem.(dockerRes).GetId()
						m.activeDialog = getRemoveVolumeDialog(map[string]string{"ID": volumeId})
						m.showDialog = true
						cmds = append(cmds, m.activeDialog.Init())
					}

				case key.Matches(assertedMsg, m.keymap.volume.DeleteForce):
					selectedItems := m.getSelectedItems()

					op := volumeDeleteBulk(m.dockerClient, selectedItems, true, m.activeTab, m.notificationChan, m.possibleLongRunningOpErrorChan)
					go m.runBackground(op)

					cmds = append(cmds, clearSelectionCmd())

				case key.Matches(assertedMsg, m.keymap.volume.CopyId):
					currentItem := m.getSelectedItem()

					if currentItem != nil && !m.isCurrentTabInBulkMode() {
						dres := currentItem.(dockerRes)
						op := copyIdToClipboard(dres, m.activeTab, m.notificationChan)
						op()
					}
				}
			}

		}

	case teadialog.CloseDialog:
		m.showDialog = false
		dialog, ok := m.activeDialog.(teadialog.Dialog)

		// if the m.active dialog is not a `teadialog.Dialog` (i.e could be `teadialog.ErrorDialog`), then do not proceed forward.
		if !ok {
			break
		}

		dialogRes := dialog.GetUserChoices()

		switch dialogRes.Kind {
		case dialogRemoveContainer:
			userChoice := dialogRes.UserChoices

			opts := it.ContainerRemoveOpts{
				RemoveVolumes: userChoice["remVols"].(bool),
				RemoveLinks:   userChoice["remLinks"].(bool),
				Force:         userChoice["force"].(bool),
			}

			containerId := dialogRes.UserStorage["ID"]
			if containerId != "" {
				op := containerDelete(m.dockerClient, containerId, opts, m.activeTab, m.notificationChan)

				go m.runBackground(op)
			}

		case dialogPruneContainers:
			userChoice := dialogRes.UserChoices

			if userChoice["confirm"] == "Yes" {
				// prune containers on a separate goroutine, since UI gets stuck otherwise(since this may take sometime)
				op := func() error {
					report, err := m.dockerClient.PruneContainers()

					if err != nil {
						return err
					}

					// we send notification directly from go routine since the main goroutine does not have access to `report`
					msg := fmt.Sprintf("Pruned %d containers", report.ContainersDeleted)
					m.notificationChan <- NewNotification(m.activeTab, listStatusMessageStyle.Render(msg))

					return nil
				}

				go m.runBackground(op)
			}

		case dialogRunImage:
			userChoices := dialogRes.UserChoices
			storage := dialogRes.UserStorage

			portMappingStr := userChoices["port"]
			portMappings, err := GetPortMappingFromStr(portMappingStr.(string))

			if err != nil {
				m.activeDialog = teadialog.NewErrorDialog(err.Error(), m.width)
				m.showDialog = true
				break
			}

			var envVars []string
			if userChoices["env"].(string) != "" {
				envVars = strings.Split(userChoices["env"].(string), ",")
			}

			config := it.ContainerCreateConfig{
				Env:          envVars,
				ImageId:      storage["ID"],
				PortBindings: portMappings,
				Name:         userChoices["name"].(string),
			}

			op := runImage(m.dockerClient, config, m.activeTab, m.notificationChan)

			go m.runBackground(op)

		case dialogPruneImages:
			userChoice := dialogRes.UserChoices

			if userChoice["confirm"] == "Yes" {
				op := imagePrune(m.dockerClient, m.activeTab, m.notificationChan)
				go m.runBackground(op)
			}

		case dialogPruneVolumes:
			userChoice := dialogRes.UserChoices

			if userChoice["confirm"] == "Yes" {
				// same reason as above, again
				op := volumePrune(m.dockerClient, m.activeTab, m.notificationChan)
				go m.runBackground(op)
			}

		case dialogRemoveVolumes:
			userChoice := dialogRes.UserChoices
			volumeId := dialogRes.UserStorage["ID"]

			if volumeId != "" {
				op := volumeDelete(m.dockerClient, volumeId, userChoice["force"].(bool), m.activeTab, m.notificationChan)
				go m.runBackground(op)
			}

		case dialogRemoveImage:
			userChoice := dialogRes.UserChoices

			imageId := dialogRes.UserStorage["ID"]

			if imageId != "" {
				opts := it.RemoveImageOptions{
					Force:   userChoice["force"].(bool),
					NoPrune: !userChoice["pruneChildren"].(bool),
				}

				op := imageDelete(m.dockerClient, imageId, opts, m.activeTab, m.notificationChan)
				go m.runBackground(op)
			}

		case dialogImageBuild:
			userChoice := dialogRes.UserChoices
			tagsStr := userChoice["image_tags"].(string)
			tags := strings.Split(tagsStr, ",")

			buildContext, _ := os.Getwd()
			options := it.ImageBuildOptions{
				Tags:       tags,
				Dockerfile: "Dockerfile",
			}

			progressModel := components.NewProgressBar()
			buildInfoCard := getBuildProgress(progressModel)
			m.activeDialog = buildInfoCard
			m.showDialog = true

			cmds = append(cmds, buildInfoCard.Init())

			op := func() error {
				res, err := m.dockerClient.BuildImage(buildContext, options)

				if err != nil {
					return err
				}

				decoder := json.NewDecoder(res.Body)

				var status it.ImageBuildJSON

				for {
					if err := decoder.Decode(&status); errors.Is(err, io.EOF) {
						break
					}

					if status.Error != nil {
						return errors.New(status.Error.Message)
					}

					buildInfoCard.progressChan <- status.Stream
				}

				/*
					HACK: I add `Step 2/1 : ` becuz we do regex matching in buildProgress.Update to extract current Step
					and calculate progress bar completion, adding `2/1` will enable the progress bar to show 100% when image
					is done building
				*/

				/*
				 FIX: this doesn't get colored, mostly prolly because buildInfoCard.Update uses regex to split string into groups,
				 so ANSI color codes are lost
				*/
				buildInfoCard.progressChan <- successForeground.Render("Step 2/1 : Build Complete!")

				m.notificationChan <- NewNotification(m.activeTab, listStatusMessageStyle.Render("Build Complete!"))

				return nil
			}

			go m.runBackground(op)
		}

	}

	// var cmd tea.Cmd
	// do not pass key.msg to list if dialog is active, otherwise tui updates to navigation keys
	if !m.showDialog {
		newList, cmd := m.TabContent[m.activeTab].Update(msg)
		m.TabContent[m.activeTab] = newList.(listModel)
		// m.TabContent[m.activeTab].list, cmd = m.TabContent[m.activeTab].list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m MainModel) runBackground(op Operation) {
	if err := op(); err != nil {
		m.possibleLongRunningOpErrorChan <- err
	}
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func (m MainModel) View() string {

	if m.windowTooSmall {
		return dialogContainerStyle.Render(m.windowtoosmallModel.View())
	}

	if m.showDialog {
		return dialogContainerStyle.Render(m.activeDialog.View())
	}
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == int(m.activeTab)
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "└"
		} else if isLast && !isActive {
			border.BottomRight = "┴"
		}

		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	var row string
	row = lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	fillerStringLen := windowStyle.GetWidth() - lipgloss.Width(row)
	if fillerStringLen > 0 {
		fillerString := strings.Repeat("─", fillerStringLen+1)
		fillerString += "┐"
		row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, fillerStyle.Render(fillerString))
	}

	list := m.TabContent[m.activeTab].View()
	curItem := m.getSelectedItem()

	infobox := ""
	if curItem != nil && m.displayInfoBox {
		infobox = m.populateInfoBox(curItem)
	}

	// TODO: align info box to right edge of the window
	body_with_info := lipgloss.JoinHorizontal(lipgloss.Top, list, infobox)

	tabKeyBinds := m.getActiveTab().getKeymap()
	tabKeyBindsStr := m.helpGen.View(tabKeyBinds)

	// we do this cuz the help string is misaligned when it is of more than 2 lines, so we split and add space to align them

	AlignHelpText := func(helpText string) string {
		substrs := strings.SplitAfter(helpText, "\n")
		if len(substrs) > 1 {
			// add space prefix to second substring
			substrs[1] = "  " + substrs[1]
			helpText = strings.Join(substrs, "")
		}

		return helpText
	}

	navKeyBinds := AlignHelpText(m.navKeymap.View(m.keymap.navigation))
	tabKeyBindsStr = AlignHelpText(tabKeyBindsStr)

	help := lipgloss.JoinVertical(lipgloss.Left, "  "+navKeyBinds, "  "+tabKeyBindsStr)
	help = lipgloss.PlaceVertical(5, lipgloss.Bottom, help)
	body_with_help := lipgloss.JoinVertical(lipgloss.Top, body_with_info, help)
	body_with_info = windowStyle.Render(body_with_help)

	doc.WriteString(row)
	doc.WriteString("\n")

	doc.WriteString(body_with_info)
	return docStyle.Render(doc.String())
}

// helpers

/*
Fetches new data from the docker api and returns []dockerRes, also updates other required fields depending on the tabId passed.
Passing `wg` is optional and is add is added for the sole purpose of testing.
*/
func (m MainModel) fetchNewData(tab tabId, wg *sync.WaitGroup) []dockerRes {
	var newlist []dockerRes
	switch tab {
	case IMAGES:
		newImgs := m.dockerClient.ListImages()
		newlist = makeImageItems(newImgs)

		// update imageToName map if there are new images
		if wg != nil {
			wg.Add(1)
		}

		go func() {
			if wg != nil {
				defer wg.Done()
			}
			for _, image := range newlist {
				if _, keyExists := m.imageIdToNameMap[image.GetId()]; !keyExists {
					m.imageIdToNameMap[image.GetId()] = image.getName()
				}
			}
		}()

	case CONTAINERS:
		newContainers := m.dockerClient.ListContainers(false)
		newlist = makeContainerItems(newContainers)

		for _, newContainer := range newlist {
			id := newContainer.GetId()
			if _, ok := m.containerSizeTracker.sizeMap[id]; !ok {

				if wg != nil {
					wg.Add(1)
				}
				go func() {
					if wg != nil {
						defer wg.Done()
					}
					containerInfo, err := m.dockerClient.InspectContainer(id)

					if err != nil {
						panic(err)
					}

					updateContainerSizeMap(*containerInfo, &m.containerSizeTracker)
				}()
			}
		}

	case VOLUMES:
		// TODO: handle errors
		newVolumes, _ := m.dockerClient.ListVolumes()
		newlist = makeVolumeItem(newVolumes)
	}

	return newlist
}

// fetches new docker data and updates the list in the current tab.
func (m MainModel) updateContent(tab tabId) MainModel {
	newlist := m.fetchNewData(tab, nil)
	// m.TabContent[tab] = m.TabContent[tab].updateTab(m.dockerClient)
	listM, _ := m.TabContent[tab].Update(newlist)
	m.TabContent[tab] = listM.(listModel)
	return m
}

// Generates info box for the current list item
func (m MainModel) populateInfoBox(item list.Item) string {
	temp, _ := item.(dockerRes)
	switch m.activeTab {
	case IMAGES:
		if it, ok := temp.(imageItem); ok {
			info := populateImageInfoBox(it)
			return moreInfoStyle.Render(info)
		}

	case CONTAINERS:
		if ct, ok := temp.(containerItem); ok {
			info := populateContainerInfoBox(ct, &m.containerSizeTracker, m.imageIdToNameMap)
			return moreInfoStyle.Render(info)
		}

	case VOLUMES:
		if vt, ok := temp.(VolumeItem); ok {
			info := populateVolumeInfoBox(vt)
			return moreInfoStyle.Render(info)
		}
	}
	return ""
}

// Util
func (m *MainModel) nextTab() {
	if int(m.activeTab) == len(CONFIG_TAB_ORDERING)-1 {
		m.activeTab = 0
	} else {
		m.activeTab += 1
	}
}

func (m *MainModel) prevTab() {
	if int(m.activeTab) == 0 {
		m.activeTab = tabId(len(CONFIG_TAB_ORDERING) - 1)
	} else {
		m.activeTab -= 1
	}
}

// Gets currently active tab (type: listModel)
func (m MainModel) getActiveTab() listModel {
	return m.TabContent[m.activeTab]
}

// Gets currently active list (type: list.Model)
func (m MainModel) getActiveList() *list.Model {
	return &m.TabContent[m.activeTab].list
}

// Get list at specified index
func (m MainModel) getList(index int) *list.Model {
	if index >= len(m.TabContent) {
		panic(fmt.Sprintf("Index %d out of bounds", index))
	}
	return &m.TabContent[index].list
}

// Helper function to get current focused item in the list
func (m MainModel) getSelectedItem() list.Item {
	return m.TabContent[m.activeTab].list.SelectedItem()
}

/*
Helper function to get selected Items in the list to perform bulk operations
if no items are selected, returns the current Item the cursor is on in a list
of length 1
*/
func (m MainModel) getSelectedItems() []dockerRes {
	activeTab := m.getActiveTab()

	if activeTab.inBulkMode() {
		selectedMap := activeTab.list.GetSelected()

		vals := make([]dockerRes, len(selectedMap))

		i := 0
		for _, val := range selectedMap {
			vals[i] = val.(dockerRes)
			i++
		}

		return vals
	} else {
		return []dockerRes{activeTab.list.SelectedItem().(dockerRes)}
	}
}

// Helper function to get state of current tab (i.e bulk mode or nah)
func (m MainModel) isCurrentTabInBulkMode() bool {
	return m.getActiveTab().inBulkMode()
}

// Copies str to clipboard
func copyToClipboard(str string) {
	str = strings.TrimPrefix(str, "sha256:")[:20]
	clipboard.Write(clipboard.FmtText, []byte(str))
}

// Prepopulates container size info at the start of the program concurrently
func (m *MainModel) prepopulateContainerSizeMapConcurrently() {
	containerInfoWithSize := m.dockerClient.ListContainers(true)

	for _, info := range containerInfoWithSize {
		m.containerSizeTracker.sizeMap[info.ID] = ContainerSize{
			sizeRw: info.SizeRw,
			rootFs: info.SizeRootFs,
		}
	}
}

// Adds size info from containerInfo to containersizeTracker. Meant to be used on demand when new container gets added.
func updateContainerSizeMap(containerInfo it.InspectContainerData, containerSizeTracker *ContainerSizeManager) {
	containerSizeTracker.mu.Lock()
	containerSizeTracker.sizeMap[containerInfo.ID] = ContainerSize{
		sizeRw: containerInfo.SizeRw,
		rootFs: containerInfo.SizeRootFs,
	}
	containerSizeTracker.mu.Unlock()
}
