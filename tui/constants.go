package tui

type KeyBindings struct {
	navigationKeyBindings    navigationKeyBindings
	imageKeyBindings         imageKeyBindings
	imageKeyBindingsBulk     imageKeyBindingsBulk
	containerKeyBindings     containerKeyBindings
	containerKeyBindingsBulk containerKeyBindingsBulk
	volumeKeyBindings        volumeKeyBindings
	volumeKeyBindingsBulk    volumeKeyBindingsBulk
}

type navigationKeyBindings struct {
	Enter    []string
	Back     []string
	Quit     []string
	Select   []string
	NextTab  []string
	PrevTab  []string
	NextItem []string
	PrevItem []string
	PrevPage []string
	NextPage []string
}

type imageKeyBindings struct {
	Run         []string
	Rename      []string
	Build       []string
	Scout       []string
	Prune       []string
	Delete      []string
	DeleteForce []string
	CopyId      []string
	RunAndExec  []string
}

type imageKeyBindingsBulk struct {
	DeleteForceBulk   []string
	ExitSelectionMode []string
}

type containerKeyBindings struct {
	ToggleListAll   []string
	ToggleStartStop []string
	TogglePause     []string
	Restart         []string
	Delete          []string
	DeleteForce     []string
	Exec            []string
	Prune           []string
	CopyId          []string
	ShowLogs        []string
}

type containerKeyBindingsBulk struct {
	ToggleListAll     []string
	ToggleStartStop   []string
	TogglePause       []string
	Restart           []string
	DeleteForce       []string
	ExitSelectionMode []string
}

type volumeKeyBindings struct {
	Delete      []string
	DeleteForce []string
	Prune       []string
	CopyId      []string
}

type volumeKeyBindingsBulk struct {
	DeleteForce       []string
	ExitSelectionMode []string
}

func initKeyBindingsConstant() *KeyBindings {
	return &KeyBindings{
		navigationKeyBindings: navigationKeyBindings{
			Enter:    globalConfig.Strings("keybindings.navigation.Enter"),
			Back:     globalConfig.Strings("keybindings.navigation.Back"),
			Quit:     globalConfig.Strings("keybindings.navigation.Quit"),
			Select:   globalConfig.Strings("keybindings.navigation.Select"),
			NextTab:  globalConfig.Strings("keybindings.navigation.NextTab"),
			PrevTab:  globalConfig.Strings("keybindings.navigation.PrevTab"),
			NextItem: globalConfig.Strings("keybindings.navigation.NextItem"),
			PrevItem: globalConfig.Strings("keybindings.navigation.PrevItem"),
			PrevPage: globalConfig.Strings("keybindings.navigation.PrevPage"),
			NextPage: globalConfig.Strings("keybindings.navigation.NextPage"),
		},
		imageKeyBindings: imageKeyBindings{
			Run:         globalConfig.Strings("keybindings.image.Run"),
			Rename:      globalConfig.Strings("keybindings.image.Rename"),
			Build:       globalConfig.Strings("keybindings.image.Build"),
			Scout:       globalConfig.Strings("keybindings.image.Scout"),
			Prune:       globalConfig.Strings("keybindings.image.Prune"),
			Delete:      globalConfig.Strings("keybindings.image.Delete"),
			DeleteForce: globalConfig.Strings("keybindings.image.DeleteForce"),
			CopyId:      globalConfig.Strings("keybindings.image.CopyId"),
			RunAndExec:  globalConfig.Strings("keybindings.image.RunAndExec"),
		},
		imageKeyBindingsBulk: imageKeyBindingsBulk{
			DeleteForceBulk:   globalConfig.Strings("keybindings.image.DeleteForce"),
			ExitSelectionMode: globalConfig.Strings("keybindings.image.ExitSelectionMode"),
		},
		containerKeyBindings: containerKeyBindings{
			ToggleListAll:   globalConfig.Strings("keybindings.container.ToggleListAll"),
			ToggleStartStop: globalConfig.Strings("keybindings.container.ToggleStartStop"),
			TogglePause:     globalConfig.Strings("keybindings.container.TogglePause"),
			Restart:         globalConfig.Strings("keybindings.container.Restart"),
			Delete:          globalConfig.Strings("keybindings.container.Delete"),
			DeleteForce:     globalConfig.Strings("keybindings.container.DeleteForce"),
			Exec:            globalConfig.Strings("keybindings.container.Exec"),
			Prune:           globalConfig.Strings("keybindings.container.Prune"),
			CopyId:          globalConfig.Strings("keybindings.container.CopyId"),
			ShowLogs:        globalConfig.Strings("keybindings.container.ShowLogs"),
		},
		containerKeyBindingsBulk: containerKeyBindingsBulk{
			ToggleListAll:     globalConfig.Strings("keybindings.containerBulk.ToggleListAll"),
			ToggleStartStop:   globalConfig.Strings("keybindings.containerBulk.ToggleStartStop"),
			TogglePause:       globalConfig.Strings("keybindings.containerBulk.TogglePause"),
			Restart:           globalConfig.Strings("keybindings.containerBulk.Restart"),
			DeleteForce:       globalConfig.Strings("keybindings.containerBulk.DeleteForce"),
			ExitSelectionMode: globalConfig.Strings("keybindings.containerBulk.ExitSelectionMode"),
		},
		volumeKeyBindings: volumeKeyBindings{
			Delete:      globalConfig.Strings("keybindings.volume.Delete"),
			DeleteForce: globalConfig.Strings("keybindings.volume.DeleteForce"),
			Prune:       globalConfig.Strings("keybindings.volume.Prune"),
			CopyId:      globalConfig.Strings("keybindings.volume.CopyId"),
		},
		volumeKeyBindingsBulk: volumeKeyBindingsBulk{
			DeleteForce:       globalConfig.Strings("keybindings.volume.DeleteForce"),
			ExitSelectionMode: globalConfig.Strings("keybindings.volume.ExitSelectionMode"),
		},
	}
}
