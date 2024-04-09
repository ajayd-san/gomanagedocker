package tui

import (
	"log"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/api/types/image"
)

type item struct {
	title string
	desc  float64
}

func makeItem(title string, desc float64) item {
	return item{title, desc}
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return strconv.FormatFloat(i.desc, 'f', 2, 64) }
func (i item) FilterValue() string { return i.title }

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
	items := []list.Item{
		item{title: "1", desc: 1},
		item{title: "2", desc: 2},
		item{title: "3", desc: 3},
	}

	m := listModel{list: list.New(items, list.NewDefaultDelegate(), 100, 36)}
	m.list.SetShowTitle(false)
	return m
}

func makeItems(raw []image.Summary) []list.Item {
	listItems := make([]list.Item, len(raw))
	log.Println(raw[0])

	//INFO: only converting to gb (might want to change later to accomidate mb)
	for i, data := range raw {
		listItems[i] = makeItem(data.ID, float64(data.Size)/float64(1e+9))
	}

	return listItems
}

// Util

func (m *listModel) updateContent(content []image.Summary) {
	m.list.SetItems(makeItems(content))
}
