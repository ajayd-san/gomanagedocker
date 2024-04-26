package tui

import "github.com/charmbracelet/bubbles/key"

type navigationKeymap struct {
	Enter key.Binding
	Back  key.Binding
	Quit  key.Binding
	Next  key.Binding
	Prev  key.Binding
}

type imgKeymap struct {
	Create key.Binding
	Rename key.Binding
	Pull   key.Binding
	Delete key.Binding
}

type contKeymap struct {
	ToggleListAll   key.Binding
	ToggleStartStop key.Binding
	Delete          key.Binding
	Exec            key.Binding
}

type volKeymap struct {
	Delete key.Binding
}

var ImageKeymapm = imgKeymap{
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
	Pull: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "Pull new Image"),
	),
}

var ContainerKeymap = contKeymap{
	ToggleListAll: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "Toggle list all"),
	),
	ToggleStartStop: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Toggle Start/Stop"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Exec: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "exec"),
	),
}

var VolumeKeymap = volKeymap{
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
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

func getVolumeKeymap() []key.Binding {
	return []key.Binding{
		VolumeKeymap.Delete,
	}
}

func getImageKeymap() []key.Binding {
	return []key.Binding{
		ImageKeymapm.Delete,
		ImageKeymapm.Pull,
	}
}

func getContainerKeymap() []key.Binding {
	return []key.Binding{
		ContainerKeymap.ToggleListAll,
		ContainerKeymap.ToggleStartStop,
		ContainerKeymap.Delete,
		ContainerKeymap.Exec,
	}
}
