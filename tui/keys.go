package tui

import (
	"github.com/charmbracelet/bubbles/key"
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

func (m navigationKeymap) FullHelp() [][]key.Binding {
	allBindings := []key.Binding{m.NextItem, m.PrevItem, m.NextTab, m.PrevTab, m.PrevPage, m.NextPage, m.Quit}
	return packKeybindings(allBindings, KeymapAvailableWidth)
}

func (m navigationKeymap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (keybinds *KeyBindings) initNavigationKeys() *navigationKeymap {
	return &navigationKeymap{
		Enter: key.NewBinding(
			key.WithKeys(keybinds.navigationKeyBindings.Enter...),
			key.WithHelp(ArrayToString(keybinds.navigationKeyBindings.Enter), "select"),
		),
		Back: key.NewBinding(
			key.WithKeys(keybinds.navigationKeyBindings.Back...),
			key.WithHelp(ArrayToString(keybinds.navigationKeyBindings.Back), "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys(keybinds.navigationKeyBindings.Quit...),
			key.WithHelp(ArrayToString(keybinds.navigationKeyBindings.Quit), "quit"),
		),
		Select: key.NewBinding(
			key.WithKeys(keybinds.navigationKeyBindings.Select...),
			key.WithHelp(ArrayToString(keybinds.navigationKeyBindings.Select), "Select"),
		),
		NextTab: key.NewBinding(
			key.WithKeys(keybinds.navigationKeyBindings.NextTab...),
			key.WithHelp(ArrayToString(keybinds.navigationKeyBindings.NextTab), "next"),
		),
		PrevTab: key.NewBinding(
			key.WithKeys(keybinds.navigationKeyBindings.PrevTab...),
			key.WithHelp(ArrayToString(keybinds.navigationKeyBindings.PrevTab), "prev"),
		),
		NextItem: key.NewBinding(
			key.WithKeys(keybinds.navigationKeyBindings.NextItem...),
			key.WithHelp(ArrayToString(keybinds.navigationKeyBindings.NextItem), "next item"),
		),
		PrevItem: key.NewBinding(
			key.WithKeys(keybinds.navigationKeyBindings.PrevItem...),
			key.WithHelp(ArrayToString(keybinds.navigationKeyBindings.PrevItem), "prev item"),
		),
		PrevPage: key.NewBinding(
			key.WithKeys(keybinds.navigationKeyBindings.NextPage...),
			key.WithHelp("[", "prev page"),
		),
		NextPage: key.NewBinding(
			key.WithKeys(keybinds.navigationKeyBindings.NextPage...),
			key.WithHelp("]", "next page"),
		),
	}
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

func (keybinds *KeyBindings) initImageKeys() *imgKeymap {
	return &imgKeymap{
		Run: key.NewBinding(
			key.WithKeys(keybinds.imageKeyBindings.Run...),
			key.WithHelp(ArrayToString(keybinds.imageKeyBindings.Run), "run"),
		),
		Rename: key.NewBinding(
			key.WithKeys(keybinds.imageKeyBindings.Rename...),
			key.WithHelp(ArrayToString(keybinds.imageKeyBindings.Rename), "rename"),
		),
		Build: key.NewBinding(
			key.WithKeys(keybinds.imageKeyBindings.Build...),
			key.WithHelp(ArrayToString(keybinds.imageKeyBindings.Build), "build"),
		),
		Delete: key.NewBinding(
			key.WithKeys(keybinds.imageKeyBindings.Delete...),
			key.WithHelp(ArrayToString(keybinds.imageKeyBindings.Delete), "delete"),
		),
		DeleteForce: key.NewBinding(
			key.WithKeys(keybinds.imageKeyBindings.DeleteForce...),
			key.WithHelp(ArrayToString(keybinds.imageKeyBindings.DeleteForce), "delete (force)"),
		),
		Scout: key.NewBinding(
			key.WithKeys(keybinds.imageKeyBindings.Scout...),
			key.WithHelp(ArrayToString(keybinds.imageKeyBindings.Scout), "scout"),
		),
		Prune: key.NewBinding(
			key.WithKeys(keybinds.imageKeyBindings.Prune...),
			key.WithHelp(ArrayToString(keybinds.imageKeyBindings.Prune), "prune images"),
		),
		CopyId: key.NewBinding(
			key.WithKeys(keybinds.imageKeyBindings.CopyId...),
			key.WithHelp(ArrayToString(keybinds.imageKeyBindings.CopyId), "copy Image ID"),
		),
		RunAndExec: key.NewBinding(
			key.WithKeys(keybinds.imageKeyBindings.RunAndExec...),
			key.WithHelp(ArrayToString(keybinds.imageKeyBindings.RunAndExec), "run and exec"),
		),
	}
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

// This not required and is there only to satisfy the
func (m imgKeymap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

type imgKeymapBulk struct {
	DeleteForce       key.Binding
	ExitSelectionMode key.Binding
}

// var imageKeymapBulk =
func (keybinds *KeyBindings) initImageKeysBulk() *imgKeymapBulk {
	return &imgKeymapBulk{
		DeleteForce: key.NewBinding(
			key.WithKeys(keybinds.imageKeyBindingsBulk.DeleteForceBulk...),
			key.WithHelp(ArrayToString(keybinds.imageKeyBindingsBulk.DeleteForceBulk), "bulk delete (force)"),
		),
		ExitSelectionMode: key.NewBinding(
			key.WithKeys(keybinds.imageKeyBindingsBulk.ExitSelectionMode...),
			key.WithHelp(ArrayToString(keybinds.imageKeyBindingsBulk.ExitSelectionMode), "exit selection mode"),
		),
	}
}

// not required, only there to satify the help.KeyMap interface
func (im imgKeymapBulk) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (im imgKeymapBulk) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{im.DeleteForce}, {im.ExitSelectionMode},
	}
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

func (keybinds *KeyBindings) initContainerKeys() *contKeymap {
	return &contKeymap{
		ToggleListAll: key.NewBinding(
			key.WithKeys(keybinds.containerKeyBindings.ToggleListAll...),
			key.WithHelp(ArrayToString(keybinds.containerKeyBindings.ToggleListAll), "toggle list all"),
		),
		ToggleStartStop: key.NewBinding(
			key.WithKeys(keybinds.containerKeyBindings.ToggleStartStop...),
			key.WithHelp(ArrayToString(keybinds.containerKeyBindings.ToggleStartStop), "toggle Start/Stop"),
		),
		TogglePause: key.NewBinding(
			key.WithKeys(keybinds.containerKeyBindings.TogglePause...),
			key.WithHelp(ArrayToString(keybinds.containerKeyBindings.TogglePause), "toggle Pause/unPause"),
		),
		Restart: key.NewBinding(
			key.WithKeys(keybinds.containerKeyBindings.Restart...),
			key.WithHelp(ArrayToString(keybinds.containerKeyBindings.Restart), "restart"),
		),
		Delete: key.NewBinding(
			key.WithKeys(keybinds.containerKeyBindings.Delete...),
			key.WithHelp(ArrayToString(keybinds.containerKeyBindings.Delete), "delete"),
		),
		DeleteForce: key.NewBinding(
			key.WithKeys(keybinds.containerKeyBindings.DeleteForce...),
			key.WithHelp(ArrayToString(keybinds.containerKeyBindings.DeleteForce), "delete (force)"),
		),
		Prune: key.NewBinding(
			key.WithKeys(keybinds.containerKeyBindings.Prune...),
			key.WithHelp(ArrayToString(keybinds.containerKeyBindings.Prune), "prune"),
		),
		Exec: key.NewBinding(
			key.WithKeys(keybinds.containerKeyBindings.Exec...),
			key.WithHelp(ArrayToString(keybinds.containerKeyBindings.Exec), "exec"),
		),
		CopyId: key.NewBinding(
			key.WithKeys(keybinds.containerKeyBindings.CopyId...),
			key.WithHelp(ArrayToString(keybinds.containerKeyBindings.CopyId), "copy ID"),
		),
		ShowLogs: key.NewBinding(
			key.WithKeys(keybinds.containerKeyBindings.ShowLogs...),
			key.WithHelp(ArrayToString(keybinds.containerKeyBindings.ShowLogs), "Show Logs"),
		),
	}
}

type contKeymapBulk struct {
	ToggleListAll     key.Binding
	ToggleStartStop   key.Binding
	TogglePause       key.Binding
	Restart           key.Binding
	DeleteForce       key.Binding
	ExitSelectionMode key.Binding
}

func (co contKeymapBulk) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (co contKeymapBulk) FullHelp() [][]key.Binding {
	bindings := []key.Binding{
		co.ToggleListAll,
		co.ToggleStartStop,
		co.Restart,
		co.TogglePause,
		co.DeleteForce,
		co.ExitSelectionMode,
	}

	return packKeybindings(bindings, KeymapAvailableWidth)
}

func (keybinds *KeyBindings) initContainerKeysBulk() *contKeymapBulk {
	return &contKeymapBulk{
		ToggleListAll: key.NewBinding(
			key.WithKeys(keybinds.containerKeyBindingsBulk.ToggleListAll...),
			key.WithHelp(ArrayToString(keybinds.containerKeyBindingsBulk.ToggleListAll), "toggle list all"),
		),
		ToggleStartStop: key.NewBinding(
			key.WithKeys(keybinds.containerKeyBindingsBulk.ToggleStartStop...),
			key.WithHelp(ArrayToString(keybinds.containerKeyBindingsBulk.ToggleStartStop), "Bulk toggle Start/Stop"),
		),
		TogglePause: key.NewBinding(
			key.WithKeys(keybinds.containerKeyBindingsBulk.TogglePause...),
			key.WithHelp(ArrayToString(keybinds.containerKeyBindingsBulk.TogglePause), "Bulk toggle Pause/unPause"),
		),
		Restart: key.NewBinding(
			key.WithKeys(keybinds.containerKeyBindingsBulk.Restart...),
			key.WithHelp(ArrayToString(keybinds.containerKeyBindingsBulk.Restart), "Bulk restart"),
		),
		DeleteForce: key.NewBinding(
			key.WithKeys(keybinds.containerKeyBindingsBulk.DeleteForce...),
			key.WithHelp(ArrayToString(keybinds.containerKeyBindingsBulk.DeleteForce), "Bulk delete (force)"),
		),
		ExitSelectionMode: key.NewBinding(
			key.WithKeys(keybinds.containerKeyBindingsBulk.ExitSelectionMode...),
			key.WithHelp(ArrayToString(keybinds.containerKeyBindingsBulk.ExitSelectionMode), "exit selection mode"),
		),
	}
}

type volKeymap struct {
	Delete      key.Binding
	DeleteForce key.Binding
	Prune       key.Binding
	CopyId      key.Binding
}

func (m volKeymap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (m volKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{m.Delete}, {m.Prune}, {m.CopyId}}
}

func (keybinds *KeyBindings) initVolumeKeys() *volKeymap {
	return &volKeymap{
		Delete: key.NewBinding(
			key.WithKeys(keybinds.volumeKeyBindings.Delete...),
			key.WithHelp(ArrayToString(keybinds.volumeKeyBindings.Delete), "delete"),
		),
		DeleteForce: key.NewBinding(
			key.WithKeys(keybinds.volumeKeyBindings.DeleteForce...),
			key.WithHelp(ArrayToString(keybinds.volumeKeyBindings.DeleteForce), "delete (force)"),
		),
		Prune: key.NewBinding(
			key.WithKeys(keybinds.volumeKeyBindings.Prune...),
			key.WithHelp(ArrayToString(keybinds.volumeKeyBindings.Prune), "prune"),
		),
		CopyId: key.NewBinding(
			key.WithKeys(keybinds.volumeKeyBindings.CopyId...),
			key.WithHelp(ArrayToString(keybinds.volumeKeyBindings.CopyId), "copy Name"),
		),
	}
}

type volKeymapBulk struct {
	DeleteForce       key.Binding
	ExitSelectionMode key.Binding
}

func (vo volKeymapBulk) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (vo volKeymapBulk) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{vo.DeleteForce}, {vo.ExitSelectionMode},
	}
}

func (keybinds *KeyBindings) initVolumeKeysBulk() *volKeymapBulk {
	return &volKeymapBulk{
		DeleteForce: key.NewBinding(
			key.WithKeys(keybinds.volumeKeyBindingsBulk.DeleteForce...),
			key.WithHelp(ArrayToString(keybinds.volumeKeyBindingsBulk.DeleteForce), "bulk delete (force)"),
		),
		ExitSelectionMode: key.NewBinding(
			key.WithKeys(keybinds.volumeKeyBindingsBulk.ExitSelectionMode...),
			key.WithHelp(ArrayToString(keybinds.volumeKeyBindingsBulk.ExitSelectionMode), "exit selection mode"),
		),
	}
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
