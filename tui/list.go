package tui

import (
	"github.com/ajayd-san/gomanagedocker/dockercmd"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type listModel struct {
	list list.Model
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := listDocStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m listModel) View() string {

	return listDocStyle.Render(m.list.View())
}

func InitList() listModel {

	items := make([]list.Item, 0)
	m := listModel{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.SetShowTitle(false)
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
		newContainers := dockerClient.ListContainers()
		newlist = makeContainerItems(newContainers)
	case volumes:
		//TODO: handle errors
		newVolumes, _ := dockerClient.ListVolumes()
		newlist = makeVolumeItem(newVolumes)
	}

	if len(m.list.Items()) != len(newlist) {
		m.list.SetItems(makeItems(newlist))
	}

	return m
}
