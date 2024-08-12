package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var KeymapAvailableWidth int

type navigationKeymap struct {
	Enter    key.Binding
	Back     key.Binding
	Quit     key.Binding
	Select   key.Binding
	NextTab  key.Binding
	PrevTab  key.Binding
	NextItem key.Binding
	PrevItem key.Binding
	PrevPage key.Binding
	NextPage key.Binding
}

type imgKeymap struct {
	Run         key.Binding
	Rename      key.Binding
	Build       key.Binding
	Scout       key.Binding
	Prune       key.Binding
	Delete      key.Binding
	DeleteForce key.Binding
	CopyId      key.Binding
	RunAndExec  key.Binding
}

type contKeymap struct {
	ToggleListAll   key.Binding
	ToggleStartStop key.Binding
	TogglePause     key.Binding
	Restart         key.Binding
	Delete          key.Binding
	DeleteForce     key.Binding
	Exec            key.Binding
	Prune           key.Binding
	CopyId          key.Binding
	ShowLogs        key.Binding
}

type volKeymap struct {
	Delete key.Binding
	Prune  key.Binding
	CopyId key.Binding
}

var ImageKeymap = imgKeymap{
	Run: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "run"),
	),
	Rename: key.NewBinding(
		key.WithKeys("R"),
		key.WithHelp("r", "rename"),
	),
	Build: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "build"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	DeleteForce: key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "delete (force)"),
	),
	Scout: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "scout"),
	),
	Prune: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "prune images"),
	),
	CopyId: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy Image ID"),
	),

	RunAndExec: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "run and exec"),
	),
}

func (m imgKeymap) FullHelp() [][]key.Binding {
	allBindings := []key.Binding{
		m.Run,
		m.Build,
		m.Delete,
		m.DeleteForce,
		m.Prune,
		m.Scout,
		m.CopyId,
		m.RunAndExec,
	}

	return packKeybindings(allBindings, KeymapAvailableWidth)

}

func (m imgKeymap) ShortHelp() []key.Binding {
	return []key.Binding{
		m.Run,
		m.Delete,
		m.DeleteForce,
		m.Prune,
		m.Scout,
		m.CopyId,
		m.RunAndExec,
	}

}

var ContainerKeymap = contKeymap{
	ToggleListAll: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "toggle list all"),
	),
	ToggleStartStop: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "toggle Start/Stop"),
	),
	TogglePause: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "toggle Pause/unPause"),
	),
	Restart: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "restart"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	DeleteForce: key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "delete (force)"),
	),
	Prune: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "prune"),
	),
	Exec: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "exec"),
	),
	CopyId: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy ID"),
	),
	ShowLogs: key.NewBinding(
		key.WithKeys("L"),
		key.WithHelp("L", "Show Logs"),
	),
}

func (m contKeymap) FullHelp() [][]key.Binding {
	bindings := []key.Binding{
		m.ToggleListAll,
		m.ToggleStartStop,
		m.Restart,
		m.TogglePause,
		m.Delete,
		m.DeleteForce,
		m.Prune,
		m.Exec,
		m.CopyId,
		m.ShowLogs,
	}

	return packKeybindings(bindings, KeymapAvailableWidth)
}

func (m contKeymap) ShortHelp() []key.Binding {
	return []key.Binding{
		m.ToggleListAll,
		m.ToggleStartStop,
		m.Restart,
		m.TogglePause,
		m.Delete,
		m.DeleteForce,
		m.Prune,
		m.Exec,
		m.CopyId,
		m.ShowLogs,
	}
}

var VolumeKeymap = volKeymap{
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Prune: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "prune"),
	),
	CopyId: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy Name"),
	),
}

func (m volKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{m.Delete}, {m.Prune}, {m.CopyId}}
}

func (m volKeymap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

var NavKeymap = navigationKeymap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("q", "quit"),
	),
	Select: key.NewBinding(
		key.WithKeys(tea.KeySpace.String()),
		key.WithHelp("<space>", "Select"),
	),
	NextTab: key.NewBinding(
		key.WithKeys("right", "l", "tab"),
		key.WithHelp("->/l/tab", "next"),
	),
	PrevTab: key.NewBinding(
		key.WithKeys("left", "h", "shift+tab"),
		key.WithHelp("<-/h/S-tab", "prev"),
	),
	NextItem: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "next item"),
	),
	PrevItem: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/h", "prev item"),
	),
	PrevPage: key.NewBinding(
		key.WithKeys("["),
		key.WithHelp("[", "prev page"),
	),
	NextPage: key.NewBinding(
		key.WithKeys("]"),
		key.WithHelp("]", "next page"),
	),
}

func (m navigationKeymap) FullHelp() [][]key.Binding {
	allBindings := []key.Binding{m.NextItem, m.PrevItem, m.NextTab, m.PrevTab, m.PrevPage, m.NextPage, m.Quit}
	return packKeybindings(allBindings, KeymapAvailableWidth)
}

func (m navigationKeymap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func packKeybindings(keybindings []key.Binding, width int) [][]key.Binding {
	res := make([][]key.Binding, len(keybindings))

	i := 0
	curWidth := width
	for _, binding := range keybindings {
		if curWidth < 20 {
			i = 0
			curWidth = width
		}

		res[i] = append(res[i], binding)
		curWidth -= len(binding.Help().Desc) + len(binding.Help().Key) + 3
		i += 1
	}

	return res
}
