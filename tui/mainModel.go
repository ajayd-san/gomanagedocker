package tui

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/ajayd-san/gomanagedocker/dockercmd"
	teadialog "github.com/ajayd-san/teaDialog"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"golang.design/x/clipboard"
)

type tabId int
type TickMsg time.Time
type preloadObjects int
type preloadSizeMap struct{}

type ContainerSize struct {
	sizeRw int64
	rootFs int64
}

type ContainerSizeManager struct {
	sizeMap map[string]ContainerSize
	mu      *sync.Mutex
}

// INFO: holds container size info that is calculated on demand

type MainModel struct {
	dockerClient        dockercmd.DockerClient
	Tabs                []string
	TabContent          []listModel
	activeTab           tabId
	width               int
	height              int
	windowTooSmall      bool
	windowtoosmallModel WindowTooSmallModel
	navKeymap           help.Model
	helpGen             help.Model
	showDialog          bool
	activeDialog        tea.Model
	// we use this to cancel dialog ops when we exit from them
	dialogOpCancel context.CancelFunc
	// this maintains a map of container image sizes
	containerSizeTracker ContainerSizeManager
	// this maps imageIds to imageNames (for legibility)
	imageIdToNameMap map[string]string
	// we use this error channel to report error for possibly long running tasks, like pruning
	possibleLongRunningOpErrorChan chan error
}

// this ticker enables us to update Docker lists items every 500ms (unless set to different value in config)
func doUpdateObjectsTick() tea.Cmd {
	return tea.Tick(CONFIG_POLLING_TIME*time.Millisecond, func(t time.Time) tea.Msg { return TickMsg(t) })
}

func (m MainModel) Init() tea.Cmd {
	// check if Docker is alive, if not, exit
	err := m.dockerClient.PingDocker()
	if err != nil {
		fmt.Printf("Error connecting to Docker daemon.\nInfo: %s\n", err.Error())
		os.Exit(1)
	}
	// initialize clipboard
	err = clipboard.Init()
	log.Println(err)
	// this command enables loading tab contents a head of time, so there is no load time while switching tabs
	preloadCmd := func() tea.Msg { return preloadObjects(0) }
	return tea.Batch(preloadCmd, doUpdateObjectsTick())
}

func NewModel() MainModel {
	contents := make([]listModel, len(CONFIG_TAB_ORDERING))

	for tabid := range CONFIG_TAB_ORDERING {
		contents[tabid] = InitList(tabId(tabid))
	}

	firstTab := contents[0].tabKind

	helper := help.New()
	NavKeymap := help.New()
	return MainModel{
		dockerClient:                   dockercmd.NewDockerClient(),
		Tabs:                           CONFIG_TAB_ORDERING,
		TabContent:                     contents,
		windowtoosmallModel:            MakeNewWindowTooSmallModel(),
		possibleLongRunningOpErrorChan: make(chan error, 10),
		helpGen:                        helper,
		navKeymap:                      NavKeymap,
		activeTab:                      firstTab,
		containerSizeTracker: ContainerSizeManager{
			sizeMap: make(map[string]ContainerSize),
			mu:      &sync.Mutex{},
		},
		imageIdToNameMap: make(map[string]string),
	}
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd

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
		if d, ok := update.(teadialog.Dialog); ok {
			m.activeDialog = d
		}

		if d, ok := update.(InfoCardWrapperModel); ok {
			m.activeDialog = d
		}

		if msg, ok := msg.(tea.KeyMsg); ok && key.Matches(msg, NavKeymap.Enter) || key.Matches(msg, NavKeymap.Back) {
			if m.dialogOpCancel != nil {
				m.dialogOpCancel()
				// this might be required, in the future
				// m.dialogOpCancel = nil
			}
			m.showDialog = false
		}

		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
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
		// if window too small set and show windowTooSmall screen
		if msg.Height < 33 || msg.Width < 169 {
			m.windowTooSmall = true
			temp, _ := m.windowtoosmallModel.Update(msg)
			m.windowtoosmallModel = temp.(WindowTooSmallModel)
		} else {
			m.windowTooSmall = false
		}

		m.width = msg.Width
		m.height = msg.Height
		windowStyle = windowStyle.
			Width(m.width - listDocStyle.GetHorizontalFrameSize() - 2).
			Height(m.height - listDocStyle.GetVerticalFrameSize() - 3)

		dialogContainerStyle = dialogContainerStyle.Width(msg.Width).Height(msg.Height)

		m.helpGen.Width = msg.Width

		// change list dimensions when window size changes
		// TODO: change width
		for index := range m.TabContent {
			listM, _ := m.TabContent[index].Update(msg)
			m.TabContent[index] = listM.(listModel)
		}

	case tea.KeyMsg:
		if !m.getActiveList().SettingFilter() && !m.showDialog {
			switch {
			case key.Matches(msg, NavKeymap.Quit):
				return m, tea.Quit
			case key.Matches(msg, NavKeymap.NextTab):
				m.nextTab()
			case key.Matches(msg, NavKeymap.PrevTab):
				m.prevTab()
			}

			if m.activeTab == IMAGES {
				switch {
				case key.Matches(msg, ImageKeymap.Run):
					curItem := m.getSelectedItem()

					if curItem != nil {
						imageId := curItem.(dockerRes).getId()

						/*
							we run on a different go routine since it may take sometime to run an image(rare case)
							and we do not want to hang the main thread
						*/
						go func() {
							err := m.dockerClient.RunImage(imageId)
							if err != nil {
								m.possibleLongRunningOpErrorChan <- err
							}
						}()
					}
				case key.Matches(msg, ImageKeymap.Delete):
					curItem := m.getSelectedItem()
					if curItem != nil {
						imageId := curItem.(dockerRes).getId()
						storage := map[string]string{"ID": imageId}
						m.activeDialog = getRemoveImageDialog(storage)
						m.showDialog = true

						cmds = append(cmds, m.activeDialog.Init())
					}

				case key.Matches(msg, ImageKeymap.DeleteForce):
					curItem := m.getSelectedItem()

					if curItem != nil {
						containerId := curItem.(dockerRes).getId()

						if containerId != "" {
							err := m.dockerClient.DeleteImage(containerId, image.RemoveOptions{
								Force:         true,
								PruneChildren: false,
							})

							if err != nil {
								m.activeDialog = teadialog.NewErrorDialog(err.Error(), m.width)
								m.showDialog = true
							}
						}
					}

				case key.Matches(msg, ImageKeymap.Prune):
					m.activeDialog = getPruneImagesDialog(make(map[string]string))
					m.showDialog = true
					cmds = append(cmds, m.activeDialog.Init())

				case key.Matches(msg, ImageKeymap.Scout):
					curItem := m.getSelectedItem()
					if curItem != nil {
						dockerRes := curItem.(dockerRes)
						imageInfo := dockerRes.(imageItem)
						imageName := imageInfo.RepoTags[0]

						ctx, cancel := context.WithCancel(context.Background())
						m.dialogOpCancel = cancel

						f := func() (*dockercmd.ScoutData, error) {

							scoutData, err := m.dockerClient.ScoutImage(ctx, imageName)

							if err != nil {
								m.possibleLongRunningOpErrorChan <- err
							}

							return scoutData, err
						}

						m.activeDialog = getImageScoutDialog(f)
						m.showDialog = true
						cmds = append(cmds, m.activeDialog.Init())
					}
				case key.Matches(msg, ImageKeymap.CopyId):
					currentItem := m.getSelectedItem()

					if currentItem != nil {
						dres := currentItem.(dockerRes)
						id := dres.getId()
						id = strings.TrimPrefix(id, "sha256:")
						clipboard.Write(clipboard.FmtText, []byte(id))
					}
				}

			} else if m.activeTab == CONTAINERS {
				switch {
				case key.Matches(msg, ContainerKeymap.ToggleListAll):
					m.dockerClient.ToggleContainerListAll()

				case key.Matches(msg, ContainerKeymap.ToggleStartStop):
					log.Println("s pressed")
					curItem := m.getSelectedItem()
					if curItem != nil {
						containerId := curItem.(dockerRes).getId()
						err := m.dockerClient.ToggleStartStopContainer(containerId)

						if err != nil {
							m.activeDialog = teadialog.NewErrorDialog(err.Error(), m.width)
							m.showDialog = true
						}

					}
				case key.Matches(msg, ContainerKeymap.TogglePause):
					curItem := m.getSelectedItem()
					if curItem != nil {

						containerId := curItem.(dockerRes).getId()
						err := m.dockerClient.TogglePauseResume(containerId)

						if err != nil {
							m.activeDialog = teadialog.NewErrorDialog(err.Error(), m.width)
							m.showDialog = true
						}
					}
				case key.Matches(msg, ContainerKeymap.Restart):
					curItem := m.getSelectedItem()
					if curItem != nil {
						log.Println("in restart")
						containerId := curItem.(dockerRes).getId()
						err := m.dockerClient.RestartContainer(containerId)

						if err != nil {
							m.activeDialog = teadialog.NewErrorDialog(err.Error(), m.width)
							m.showDialog = true
						}
					}
				case key.Matches(msg, ContainerKeymap.Delete):
					curItem := m.getSelectedItem()
					if containerInfo, ok := curItem.(dockerRes); ok {
						dialog := getRemoveContainerDialog(map[string]string{"ID": containerInfo.getId()})
						m.activeDialog = dialog
						m.showDialog = true
						cmds = append(cmds, m.activeDialog.Init())
					}

				case key.Matches(msg, ContainerKeymap.DeleteForce):
					curItem := m.getSelectedItem()
					if containerInfo, ok := curItem.(dockerRes); ok {
						err := m.dockerClient.DeleteContainer(containerInfo.getId(), container.RemoveOptions{
							RemoveVolumes: false,
							RemoveLinks:   false,
							Force:         true,
						})

						if err != nil {
							m.activeDialog = teadialog.NewErrorDialog(err.Error(), m.width)
							m.showDialog = true
						}
					}

				case key.Matches(msg, ContainerKeymap.Prune):
					m.activeDialog = getPruneContainersDialog(make(map[string]string))
					m.showDialog = true
					cmds = append(cmds, m.activeDialog.Init())

				case key.Matches(msg, ContainerKeymap.Exec):
					curItem := m.getSelectedItem()
					if curItem != nil {
						containerId := curItem.(dockerRes).getId()
						// execs into the default shell of the container (got from lazydocker)
						cmd := exec.Command("docker", "exec", "-it", containerId, "/bin/sh", "-c", "eval $(grep ^$(id -un): /etc/passwd | cut -d : -f 7-)")
						cmds = append(cmds, tea.ExecProcess(cmd, nil))
					}

				case key.Matches(msg, ContainerKeymap.CopyId):
					currentItem := m.getSelectedItem()

					if currentItem != nil {
						dres := currentItem.(dockerRes)
						id := dres.getId()
						clipboard.Write(clipboard.FmtText, []byte(id))
					}
				}

			} else if m.activeTab == VOLUMES {
				switch {
				case key.Matches(msg, VolumeKeymap.Prune):
					log.Println("Volume prune called")
					curItem := m.getSelectedItem()
					if curItem != nil {
						volumeId := curItem.(dockerRes).getId()
						m.activeDialog = getPruneVolumesDialog(map[string]string{"ID": volumeId})
						m.showDialog = true
					}

				case key.Matches(msg, VolumeKeymap.Delete):
					log.Println("volume delete called")

					curItem := m.getSelectedItem()

					if curItem != nil {
						volumeId := curItem.(dockerRes).getId()
						m.activeDialog = getRemoveVolumeDialog(map[string]string{"ID": volumeId})
						m.showDialog = true
						cmds = append(cmds, m.activeDialog.Init())
					}

				case key.Matches(msg, VolumeKeymap.CopyId):
					currentItem := m.getSelectedItem()

					if currentItem != nil {
						dres := currentItem.(dockerRes)
						name := dres.getId()
						clipboard.Write(clipboard.FmtText, []byte(name))
					}
				}
			}

		}
	case teadialog.DialogSelectionResult:
		dialogRes := msg
		switch dialogRes.Kind {
		case dialogRemoveContainer:
			log.Println("remove container instruction received")
			userChoice := dialogRes.UserChoices

			opts := container.RemoveOptions{
				RemoveVolumes: userChoice["remVols"].(bool),
				RemoveLinks:   userChoice["remLinks"].(bool),
				Force:         userChoice["force"].(bool),
			}

			containerId := dialogRes.UserStorage["ID"]
			if containerId != "" {
				log.Println("removing container: ", dialogRes.UserStorage["ID"])
				err := m.dockerClient.DeleteContainer(containerId, opts)
				log.Println("container delete")
				if err != nil {
					m.activeDialog = teadialog.NewErrorDialog(err.Error(), m.width)
					m.showDialog = true
				}
			}

		case dialogPruneContainers:
			log.Println("prune containers called")
			userChoice := dialogRes.UserChoices

			if userChoice["confirm"] == "Yes" {
				log.Println("prune containers confirmed")

				// prune containers on a separate goroutine, since UI gets stuck otherwise(since this may take sometime)
				go func() {
					report, err := m.dockerClient.PruneContainers()

					log.Println(report)

					if err != nil {
						m.possibleLongRunningOpErrorChan <- err
					}
				}()
			}

		case dialogPruneImages:
			log.Println("prune images called")

			userChoice := dialogRes.UserChoices

			if userChoice["confirm"] == "Yes" {
				// run on a different go routine, same reason as above (for Prune containers)
				go func() {
					report, err := m.dockerClient.PruneImages()

					// TODO: show report on screen
					log.Println("prune images report", report)

					if err != nil {
						m.possibleLongRunningOpErrorChan <- err
					}
				}()
			}

		case dialogPruneVolumes:
			log.Println("prune volumes called")

			userChoice := dialogRes.UserChoices

			if userChoice["confirm"] == "Yes" {
				// same reason as above, again
				go func() {
					report, err := m.dockerClient.PruneVolumes()

					log.Println(report)

					if err != nil {
						m.possibleLongRunningOpErrorChan <- err
					}
				}()
			}

		case dialogRemoveVolumes:
			log.Println("remove volume called 2")
			userChoice := dialogRes.UserChoices

			volumeId := dialogRes.UserStorage["ID"]

			if volumeId != "" {
				err := m.dockerClient.DeleteVolume(volumeId, userChoice["force"].(bool))

				if err != nil {
					m.activeDialog = teadialog.NewErrorDialog(err.Error(), m.width)
					m.showDialog = true
				}
			}

		case dialogRemoveImage:
			log.Println("remove image instruction received")
			userChoice := dialogRes.UserChoices

			imageId := dialogRes.UserStorage["ID"]

			if imageId != "" {
				opts := image.RemoveOptions{
					Force:         userChoice["force"].(bool),
					PruneChildren: userChoice["pruneChildren"].(bool),
				}

				err := m.dockerClient.DeleteImage(imageId, opts)
				if err != nil {
					m.activeDialog = teadialog.NewErrorDialog(err.Error(), m.width)
					m.showDialog = true
				}
			}
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
	if curItem != nil {
		infobox = m.populateInfoBox(curItem)
		infobox = moreInfoStyle.Render(infobox)
	}

	// TODO: align info box to right edge of the window
	body_with_info := lipgloss.JoinHorizontal(lipgloss.Top, list, infobox)

	tabSpecificKeyBinds := ""

	switch m.activeTab {
	case IMAGES:
		tabSpecificKeyBinds = m.helpGen.View(ImageKeymap)
	case CONTAINERS:
		tabSpecificKeyBinds = m.helpGen.View(ContainerKeymap)
	case VOLUMES:
		tabSpecificKeyBinds = m.helpGen.View(VolumeKeymap)
	}

	body_with_help := lipgloss.JoinVertical(lipgloss.Top, body_with_info, "  "+m.navKeymap.View(NavKeymap), "  "+tabSpecificKeyBinds)
	body_with_info = windowStyle.Render(body_with_help)

	doc.WriteString(row)
	doc.WriteString("\n")

	doc.WriteString(body_with_info)
	return docStyle.Render(doc.String())
}

// helpers
func (m MainModel) fetchNewData(tab tabId) []dockerRes {
	var newlist []dockerRes
	switch tab {
	case IMAGES:
		newImgs := m.dockerClient.ListImages()
		newlist = makeImageItems(newImgs)

		// update imageToName map if there are new images
		go func() {
			for _, image := range newlist {
				if _, keyExists := m.imageIdToNameMap[image.getId()]; !keyExists {
					m.imageIdToNameMap[image.getId()] = image.getName()
				}
			}
		}()
	case CONTAINERS:
		newContainers := m.dockerClient.ListContainers(false)
		newlist = makeContainerItems(newContainers)

		for _, newContainer := range newlist {
			id := newContainer.getId()
			if _, ok := m.TabContent[CONTAINERS].ExistingIds[id]; !ok {
				go func() {
					containerInfo, err := m.dockerClient.InspectContainer(id)

					if err != nil {
						panic(err)
					}

					updateContainerSizeMap(containerInfo, &m.containerSizeTracker)
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

func (m MainModel) updateContent(tab tabId) MainModel {
	newlist := m.fetchNewData(tab)
	// m.TabContent[tab] = m.TabContent[tab].updateTab(m.dockerClient)
	listM, _ := m.TabContent[tab].Update(newlist)
	m.TabContent[tab] = listM.(listModel)
	return m
}

func (m MainModel) populateInfoBox(item list.Item) string {
	temp, _ := item.(dockerRes)
	switch m.activeTab {
	case IMAGES:
		if it, ok := temp.(imageItem); ok {
			return populateImageInfoBox(it)
		}

	case CONTAINERS:
		if ct, ok := temp.(containerItem); ok {
			return populateContainerInfoBox(ct, &m.containerSizeTracker, m.imageIdToNameMap)
		}

	case VOLUMES:
		if vt, ok := temp.(VolumeItem); ok {
			return populateVolumeInfoBox(vt)
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

func (m MainModel) getActiveTab() listModel {
	return m.TabContent[m.activeTab]
}

func (m MainModel) getActiveList() *list.Model {
	return &m.TabContent[m.activeTab].list
}

func (m MainModel) getList(index int) *list.Model {
	if index >= len(m.TabContent) {
		panic(fmt.Sprintf("Index %d out of bounds", index))
	}
	return &m.TabContent[index].list
}

func (m MainModel) getSelectedItem() list.Item {
	return m.TabContent[m.activeTab].list.SelectedItem()
}

func (m *MainModel) prepopulateContainerSizeMapConcurrently() {
	containerInfoWithSize := m.dockerClient.ListContainers(true)

	for _, info := range containerInfoWithSize {
		m.containerSizeTracker.sizeMap[info.ID] = ContainerSize{
			sizeRw: info.SizeRw,
			rootFs: info.SizeRootFs,
		}
	}
}

func updateContainerSizeMap(containerInfo *types.ContainerJSON, containerSizeTracker *ContainerSizeManager) {
	containerSizeTracker.mu.Lock()
	containerSizeTracker.sizeMap[containerInfo.ID] = ContainerSize{
		sizeRw: *containerInfo.SizeRw,
		rootFs: *containerInfo.SizeRootFs,
	}
	containerSizeTracker.mu.Unlock()
}
