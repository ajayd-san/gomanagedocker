package tui

import (
	"log"
	"slices"

	"github.com/ajayd-san/gomanagedocker/tui/components/list"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

const listWidthRatioWithInfoBox = 0.3
const listWidthRatioWithOutInfoBox = 0.85

var (
	// list always takes up 30% of the screen by default
	listWidthRatio float32 = listWidthRatioWithInfoBox
)

type itemSelect struct{}
type clearSelection struct{}

type listModel struct {
	list           list.Model
	ExistingIds    map[string]struct{}
	tabKind        tabId
	listEmpty      bool
	objectHelp     help.KeyMap
	objectHelpBulk help.KeyMap
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(int(listWidthRatio*float32(msg.Width)), msg.Height-12)
		listContainer = listContainer.Width(int(listWidthRatio * float32(msg.Width))).Height(msg.Height - 12)
	case []dockerRes:
		m.updateTab(msg)

		if len(msg) == 0 {
			m.listEmpty = true
		} else {
			m.listEmpty = false
		}
	case itemSelect:
		m.list.ToggleSelect()

		for k, _ := range m.list.GetSelected() {
			log.Println(k)
		}
		log.Println("---------------------")
		log.Println(len(m.list.GetSelected()))

	case clearSelection:
		m.list.ClearSelection()
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m listModel) View() string {
	if m.listEmpty {
		return listContainer.Render(emptyListStyle.Render("No items"))
	}

	return listContainer.Render(listDocStyle.Render(m.list.View()))
}

func (m listModel) inSelectionMode() bool {
	return len(m.list.GetSelected()) > 1
}

func InitList(tabkind tabId, objectHelp, objectHelpBulk help.KeyMap) listModel {

	items := make([]list.Item, 0)
	m := listModel{
		list:           list.New(items, list.NewDefaultDelegate(), 60, 30),
		ExistingIds:    make(map[string]struct{}),
		tabKind:        tabkind,
		objectHelp:     objectHelp,
		objectHelpBulk: objectHelpBulk,
	}

	m.list.Title = CONFIG_POLLING_TIME.String()
	m.list.StatusMessageLifetime = CONFIG_NOTIFICATION_TIMEOUT
	m.list.DisableQuitKeybindings()
	m.list.SetShowHelp(false)
	m.list.KeyMap.NextPage = key.NewBinding(key.WithKeys("]"))
	m.list.KeyMap.PrevPage = key.NewBinding(key.WithKeys("["))

	return m
}

// returns `help.KeyMap` depending on context (i.e in bulk selection mode or nah)
func (m listModel) getKeymap() help.KeyMap {
	if m.inSelectionMode() {
		return m.objectHelpBulk
	}
	return m.objectHelp
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
This function  repopulates the tab with updated items(if they are any).
For now does a linear search if the number of items have not changed to update the list (O(n) time)
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
		m.ExistingIds[item.GetId()] = struct{}{}
	}
}
