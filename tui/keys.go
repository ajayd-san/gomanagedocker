package tui

import "github.com/charmbracelet/bubbles/key"

type navigationKeymap struct {
	Enter key.Binding
	Back  key.Binding
	Quit  key.Binding
	Next  key.Binding
	Prev  key.Binding
}

type manageKeymap struct {
	Create          key.Binding
	Rename          key.Binding
	Delete          key.Binding
	ToggleStartStop key.Binding
	Exec            key.Binding
	Pull            key.Binding
}

var ManageKeymap = manageKeymap{
	Create: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "create"),
	),
	Rename: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rename"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	ToggleStartStop: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "start"),
	),
	Exec: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "exec"),
	),
	Pull: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "Pull new Image"),
	),
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
		key.WithHelp("ctrl+c/q", "quit"),
	),
	Next: key.NewBinding(
		key.WithKeys("right", "l", "tab"),
		key.WithHelp("->/l/tab", "next"),
	),
	Prev: key.NewBinding(
		key.WithKeys("left", "h", "shift+tab"),
		key.WithHelp("<-/h/shift+tab", "prev"),
	),
}

func volumeKeymap() []key.Binding {
	return []key.Binding{
		ManageKeymap.ToggleStartStop,
		ManageKeymap.Delete,
		ManageKeymap.Exec,
	}
}

func imageKeymap() []key.Binding {
	return []key.Binding{
		ManageKeymap.Delete,
		ManageKeymap.Pull,
	}
}

func containerKeymap() []key.Binding {
	return []key.Binding{
		ManageKeymap.ToggleStartStop,
		ManageKeymap.Delete,
		ManageKeymap.Exec,
	}
}
