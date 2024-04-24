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
	Create key.Binding
	Rename key.Binding
	Delete key.Binding
	Start  key.Binding
	Stop   key.Binding
	Exec   key.Binding
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
	Start: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "start")),
	Stop: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "stop")),
	Exec: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "exec")),
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
