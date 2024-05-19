package tui

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ajayd-san/gomanagedocker/dockercmd"
	dialog "github.com/ajayd-san/teaDialog"
	teadialog "github.com/ajayd-san/teaDialog"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
)

type tabId int
type TickMsg time.Time
type preloadObjects int

const (
	images tabId = iota
	containers
	volumes
)

type Model struct {
	dockerClient dockercmd.DockerClient
	Tabs         []string
	TabContent   []listModel
	activeTab    int
	width        int
	height       int
	showDialog   bool
	activeDialog tea.Model
}

func doUpdateObjectsTick() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg { return TickMsg(t) })
}

func (m Model) Init() tea.Cmd {
	preloadCmd := func() tea.Msg { return preloadObjects(0) }
	return tea.Batch(preloadCmd, doUpdateObjectsTick())
}

func NewModel(tabs []string) Model {
	contents := make([]listModel, 3)

	for i, tabKind := range []tabId{images, containers, volumes} {
		contents[i] = InitList(tabKind)
	}

	return Model{
		dockerClient: dockercmd.NewDockerClient(),
		Tabs:         tabs,
		TabContent:   contents,
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd
	//INFO: if m.showDialog is true, then hijack all keyinputs and forward them to the dialog
	//BUG: the dialog box remains fixed after first call to dialog
	if m.showDialog {
		update, cmd := m.activeDialog.Update(msg)
		if d, ok := update.(dialog.Dialog); ok {
			m.activeDialog = d
		}

		if msg, ok := msg.(tea.KeyMsg); ok && key.Matches(msg, NavKeymap.Enter) || key.Matches(msg, NavKeymap.Back) {
			m.showDialog = false
			// return m, nil
		}
		// return m, cmd

		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	//preloads all tabs, so no delay in displaying objects when first changing tabs
	case preloadObjects:
		m = m.updateContent(0)
		m = m.updateContent(1)
		m = m.updateContent(2)

	case TickMsg:
		m = m.updateContent(m.activeTab)

		// return m, doUpdateObjectsTick()
		cmds = append(cmds, doUpdateObjectsTick())

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		windowStyle = windowStyle.
			Width(m.width - listDocStyle.GetHorizontalFrameSize() - 2).
			Height(m.height - listDocStyle.GetVerticalFrameSize() - 3)

		dialogContainerStyle = dialogContainerStyle.Width(msg.Width).Height(msg.Height)

		//change list dimentions when window size changes
		// TODO: change width
		for index := range m.TabContent {
			m.getList(index).SetWidth(msg.Width)
			m.getList(index).SetHeight(msg.Height - 7)
		}

		// return m, nil

	case tea.KeyMsg:
		//INFO: if m.showDialog is true, then hijack all keyinputs and forward them to the dialog
		// if m.showDialog {
		// 	update, cmd := m.activeDialog.Update(msg)
		// 	if d, ok := update.(dialog.Dialog); ok {
		// 		m.activeDialog = d
		// 	}

		// 	if key.Matches(msg, NavKeymap.Enter) || key.Matches(msg, NavKeymap.Back) {
		// 		m.showDialog = false
		// 	}
		// 	return m, cmd
		// }
		if !m.getActiveList().SettingFilter() && !m.showDialog {
			switch {
			case key.Matches(msg, NavKeymap.Quit):
				log.Println(msg, "quitting")

				return m, tea.Quit
			case key.Matches(msg, NavKeymap.Next):
				m.nextTab()
			case key.Matches(msg, NavKeymap.Prev):
				m.prevTab()
			}

			if m.activeTab == int(images) {
				switch {
				case key.Matches(msg, ImageKeymap.Delete):
					curItem := m.getSelectedItem()
					if curItem != nil {
						imageId := curItem.(dockerRes).getId()
						storage := map[string]string{"ID": imageId}
						m.activeDialog = getRemoveImageDialog(storage)
						m.showDialog = true

						// return m, m.activeDialog.Init()
						cmds = append(cmds, m.activeDialog.Init())
					}

				case key.Matches(msg, ImageKeymap.DeleteForce):
					curItem := m.getSelectedItem()
					containerId := curItem.(dockerRes).getId()

					if containerId != "" {
						err := m.dockerClient.DeleteImage(containerId, image.RemoveOptions{
							Force:         true,
							PruneChildren: false,
						})

						if err != nil {
							m.activeDialog = teadialog.NewErrorDialog(err.Error())
							m.showDialog = true
						}
					}
				case key.Matches(msg, ImageKeymap.Prune):
					m.activeDialog = getPruneImagesDialog(make(map[string]string))
					m.showDialog = true
					// return m, m.activeDialog.Init()
					cmds = append(cmds, m.activeDialog.Init())
				}

			} else if m.activeTab == int(containers) {
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
							m.activeDialog = teadialog.NewErrorDialog(err.Error())
							m.showDialog = true
						}

					}
				case key.Matches(msg, ContainerKeymap.Delete):
					//BUG: check if any items are displayed, crashes otherwise
					curItem := m.getSelectedItem()
					containerId := curItem.(dockerRes).getId()
					dialog := getRemoveContainerDialog(map[string]string{"ID": containerId})
					m.activeDialog = dialog
					m.showDialog = true
					// return m, m.activeDialog.Init()
					cmds = append(cmds, m.activeDialog.Init())

				case key.Matches(msg, ContainerKeymap.DeleteForce):
					curItem := m.getSelectedItem()
					containerId := curItem.(dockerRes).getId()
					err := m.dockerClient.DeleteContainer(containerId, container.RemoveOptions{
						RemoveVolumes: false,
						RemoveLinks:   false,
						Force:         true,
					})

					if err != nil {
						m.activeDialog = teadialog.NewErrorDialog(err.Error())
						m.showDialog = true
					}

				case key.Matches(msg, ContainerKeymap.Prune):
					m.activeDialog = getPruneContainersDialog(make(map[string]string))
					m.showDialog = true
					// return m, m.activeDialog.Init()
					cmds = append(cmds, m.activeDialog.Init())
				}

			} else {

			}

		}
	case dialog.DialogSelectionResult:
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
				log.Println("contianer delete")
				if err != nil {
					m.activeDialog = teadialog.NewErrorDialog(err.Error())
					m.showDialog = true
				}
			}

			// return m, nil

		case dialogPruneContainers:
			log.Println("prune containers called")
			userChoice := dialogRes.UserChoices

			if userChoice["confirm"] == "Yes" {
				log.Println("prune containers confirmed")

				report, err := m.dockerClient.PruneContainers()

				log.Println(report)

				if err != nil {
					m.activeDialog = teadialog.NewErrorDialog(err.Error())
					m.showDialog = true
				}
			}

		case dialogPruneImages:
			log.Println("prune images called")

			userChoice := dialogRes.UserChoices

			if userChoice["confirm"] == "Yes" {
				report, err := m.dockerClient.PruneImages()

				//TODO: show report on screen
				log.Println("prune images report", report)

				if err != nil {
					m.activeDialog = teadialog.NewErrorDialog(err.Error())
					m.showDialog = true
				}

			}

		case dialogRemoveImage:
			log.Println("remove image instruction recieved")
			userChoice := dialogRes.UserChoices

			imageId := dialogRes.UserStorage["ID"]

			if imageId != "" {
				opts := image.RemoveOptions{
					Force:         userChoice["force"].(bool),
					PruneChildren: userChoice["pruneChildren"].(bool),
				}

				err := m.dockerClient.DeleteImage(imageId, opts)
				if err != nil {
					m.activeDialog = teadialog.NewErrorDialog(err.Error())
					m.showDialog = true
				}
			}
		}

	}

	var cmd tea.Cmd
	//do not pass key.msg to list if dialog is active
	if !m.showDialog {
		m.TabContent[m.activeTab].list, cmd = m.TabContent[m.activeTab].list.Update(msg)
	}

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func (m Model) View() string {

	if m.showDialog {
		return dialogContainerStyle.Render(m.activeDialog.View())
	}
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.activeTab
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
	infobox := PopulateInfoBox(tabId(m.activeTab), curItem)
	infobox = moreInfoStyle.Render(infobox)

	//TODO: align info box to right edge of the window
	body_with_info := lipgloss.JoinHorizontal(lipgloss.Top, list, infobox)
	// body_with_info = windowStyle.Render(body_with_info)

	doc.WriteString(row)
	doc.WriteString("\n")

	doc.WriteString(body_with_info)
	return docStyle.Render(doc.String())
}

// helpers

func (m Model) updateContent(currentTab int) Model {
	m.TabContent[currentTab] = m.TabContent[currentTab].updateTab(m.dockerClient, tabId(currentTab))
	return m
}

//Util

func (m *Model) nextTab() {
	if m.activeTab == int(volumes) {
		m.activeTab = int(images)
	} else {
		m.activeTab += 1
	}
}

func (m *Model) prevTab() {
	if m.activeTab == int(images) {
		m.activeTab = int(volumes)
	} else {
		m.activeTab -= 1
	}
}

func (m Model) getActiveTab() listModel {
	return m.TabContent[m.activeTab]
}

func (m Model) getActiveList() *list.Model {
	return &m.TabContent[m.activeTab].list
}

func (m Model) getList(index int) *list.Model {
	if index >= len(m.TabContent) {
		panic(fmt.Sprintf("Index %d out of bounds", index))
	}
	return &m.TabContent[index].list
}

func (m Model) getSelectedItem() list.Item {
	return m.TabContent[m.activeTab].list.SelectedItem()
}
