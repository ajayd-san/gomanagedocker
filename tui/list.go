package tui

import (
	"slices"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const listWidthRatio float32 = 0.3

type listModel struct {
	list        list.Model
	ExistingIds map[string]struct{}
	tabKind     tabId
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(int(listWidthRatio*float32(msg.Width)), msg.Height-10)
		listContainer = listContainer.Width(int(listWidthRatio * float32(msg.Width)))
	case []dockerRes:
		m.updateTab(msg)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m listModel) View() string {
	return listContainer.Render(listDocStyle.Render(m.list.View()))
}

func InitList(tabkind tabId) listModel {

	items := make([]list.Item, 0)
	m := listModel{
		list:        list.New(items, list.NewDefaultDelegate(), 60, 30),
		ExistingIds: make(map[string]struct{}),
		tabKind:     tabkind,
	}

	m.list.SetShowTitle(false)
	m.list.DisableQuitKeybindings()
	m.list.SetShowHelp(false)
	m.list.KeyMap.NextPage = key.NewBinding(key.WithKeys("]"))
	m.list.KeyMap.PrevPage = key.NewBinding(key.WithKeys("["))

	switch tabkind {
	case IMAGES:
		m.list.AdditionalFullHelpKeys = getImageKeymap
	case CONTAINERS:
		m.list.AdditionalFullHelpKeys = getContainerKeymap
	case VOLUMES:
		m.list.AdditionalFullHelpKeys = getVolumeKeymap
	}
	return m
}

func makeItems(raw []dockerRes) []list.Item {
	listItems := make([]list.Item, len(raw))

	// TODO: only converting to gb (might want to change later to accommodate mb)
	for i, data := range raw {
		listItems[i] = list.Item(data)
	}

	return listItems
}

// Util

/*
This function calls the docker api and repopulates the tab with updated items(if they are any).
For now does a linear search if the number of items have not changed to update the list (O(n) time)
Also, computes storage sizes for newly added containers and maps imageIds to imageNames
(to display in container infobox) in another go routine
*/
func (m *listModel) updateTab(newlist []dockerRes) {
	comparisonFunc := func(a dockerRes, b list.Item) bool {
		switch m.tabKind {
		case IMAGES:
			newA := a.(imageItem)
			newB := b.(imageItem)

			if newA.Containers != newB.Containers {
				return false
			}
		case CONTAINERS:
			newA := a.(containerItem)
			newB := b.(containerItem)

			if newA.State != newB.State {
				return false
			}
		case VOLUMES:
			// newA := a.(VolumeItem)
			// newB := b.(VolumeItem)

		}

		return true
	}

	if !slices.EqualFunc(newlist, m.list.Items(), comparisonFunc) {
		newlistItems := makeItems(newlist)
		m.list.SetItems(newlistItems)
		go m.updateExistigIds(&newlist)
	}

}

func (m *listModel) updateExistigIds(newlistItems *[]dockerRes) {
	for _, item := range *newlistItems {
		m.ExistingIds[item.getId()] = struct{}{}
	}
}
