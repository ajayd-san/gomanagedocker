package tui

import (
	"log"
	"slices"

	"github.com/ajayd-san/gomanagedocker/dockercmd"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type listModel struct {
	list        list.Model
	previousIds map[string]struct{}
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := listDocStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		// m.list.SetSize(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m listModel) View() string {
	return listDocStyle.Render(m.list.View())
}

func InitList(tab tabId) listModel {

	items := make([]list.Item, 0)
	m := listModel{list: list.New(items, list.NewDefaultDelegate(), 10, 30), previousIds: make(map[string]struct{})}

	m.list.SetShowTitle(false)
	m.list.DisableQuitKeybindings()

	switch tab {
	case images:
		m.list.AdditionalFullHelpKeys = getImageKeymap
	case containers:
		m.list.AdditionalFullHelpKeys = getContainerKeymap
	case volumes:
		m.list.AdditionalFullHelpKeys = getVolumeKeymap
	}
	return m
}

func makeItems(raw []dockerRes) []list.Item {
	listItems := make([]list.Item, len(raw))

	//TODO: only converting to gb (might want to change later to accomidate mb)
	for i, data := range raw {
		listItems[i] = list.Item(data)
	}

	return listItems
}

// Util
func (m listModel) updateTab(dockerClient dockercmd.DockerClient, id tabId) listModel {
	var newlist []dockerRes
	switch id {
	case images:
		newImgs := dockerClient.ListImages()
		newlist = makeImageItems(newImgs)
	case containers:
		newContainers := dockerClient.ListContainers(showContainerSize)
		newlist = makeContainerItems(newContainers)

		for _, newContainer := range newlist {
			id := newContainer.getId()
			if _, ok := m.previousIds[id]; !ok {
				go func() {
					containerInfo, err := dockerClient.InspectContainer(id)

					if err != nil {
						panic(err)
					}

					log.Println(containerInfo.SizeRw)
					updateContainerSizeMap(containerInfo)
				}()
			}
		}
	case volumes:
		//TODO: handle errors
		newVolumes, _ := dockerClient.ListVolumes()
		newlist = makeVolumeItem(newVolumes)
	}

	comparisionFunc := func(a dockerRes, b list.Item) bool {
		switch id {
		case images:
			newA := a.(imageItem)
			newB := b.(imageItem)

			if newA.Containers != newB.Containers {
				return false
			}
		case containers:
			newA := a.(containerItem)
			newB := b.(containerItem)

			if newA.State != newB.State {
				return false
			}
		case volumes:
			// newA := a.(VolumeItem)
			// newB := b.(VolumeItem)

		}

		return true
	}

	if !slices.EqualFunc(newlist, m.list.Items(), comparisionFunc) {
		newlistItems := makeItems(newlist)
		m.list.SetItems(newlistItems)
		go m.updateIds(newlist)
	}

	return m
}

func (m *listModel) updateIds(newlistItems []dockerRes) {
	for _, item := range newlistItems {
		m.previousIds[item.getId()] = struct{}{}
	}
}
