package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type item struct {
	title, desc, info string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
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
		item{title: "1", desc: "I have ’em all over my house"},
		item{title: "2", desc: "It's good on toast"},
		item{title: "3", desc: "It cools you down"},
		item{title: "4", desc: "And by that I mean socks without holes"},
		item{title: "5", desc: "I had this once"},
		item{title: "6", desc: "I had this once"},
		item{title: "7", desc: "Usually"},
		item{title: "8", desc: "I had this once"},
		item{title: "9", desc: "Usually"},
		item{title: "10", desc: "Usually"},
		item{title: "11 hours of sleep", desc: "I had this once"},
		item{title: "12", desc: "Usually"},
	}

	m := listModel{list: list.New(items, list.NewDefaultDelegate(), 100, 36)}
	m.list.SetShowTitle(false)
	return m
}

