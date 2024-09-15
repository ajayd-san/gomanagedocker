package tui

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	config "github.com/ajayd-san/gomanagedocker/config"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/knadh/koanf/v2"
)

const xdgPathTail string = "/gomanagedocker/gomanagedocker.yaml"

type TabOrderingMap map[string]tabId

var (
	IMAGES     tabId
	CONTAINERS tabId
	VOLUMES    tabId
)

var (
	CONFIG_POLLING_TIME         time.Duration
	CONFIG_TAB_ORDERING         []string
	CONFIG_NOTIFICATION_TIMEOUT time.Duration
	NavKeymap                   *navigationKeymap
	ImageKeymap                 *imgKeymap
	imageKeymapBulk             *imgKeymapBulk
	ContainerKeymap             *contKeymap
	ContainerKeymapBulk         *contKeymapBulk
	VolumeKeymap                *volKeymap
	volumeKeymapBulk            *volKeymapBulk
)

var globalConfig = koanf.New(".")

/*
stores fatal error that we can print before quitting gracefully
I dont think there is a native way that bubble tea lets you do it for now
*/
var earlyExitErr error

func StartTUI(debug bool) error {
	if debug {
		f, _ := tea.LogToFile("gmd_debug.log", "debug")
		defer f.Close()
	} else {
		log.SetOutput(io.Discard)
	}

	fmt.Println(tea.KeySpace.String())
	readConfig()
	loadConfig()
	loadKeyBindings()

	m := NewModel()
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		return err
	}

	/*
		we check if there is a fatal error (mostly if docker ping returned an error), print it
		and exit with non-zero error code
	*/
	if earlyExitErr != nil {
		fmt.Println(earlyExitErr.Error())
		os.Exit(1)
	}
	return nil
}

func readConfig() {
	configPath, err := os.UserConfigDir()
	if err != nil {
		log.Println("$HOME could not be determined")
	}

	config.ReadConfig(globalConfig, configPath+xdgPathTail)
}

func loadConfig() {
	CONFIG_POLLING_TIME = globalConfig.Duration("config.Polling-Time") * time.Millisecond
	CONFIG_NOTIFICATION_TIMEOUT = globalConfig.Duration("config.Notification-Timeout") * time.Millisecond
	// I have no idea how I made this work this late in the dev process, need a reliable way to test this
	CONFIG_TAB_ORDERING = globalConfig.Strings("config.Tab-Order")
	setTabConstants(CONFIG_TAB_ORDERING)
}

func loadKeyBindings() {
	keybinds := initKeyBindingsConstant()
	NavKeymap = keybinds.initNavigationKeys()
	ImageKeymap = keybinds.initImageKeys()
	imageKeymapBulk = keybinds.initImageKeysBulk()
	ContainerKeymap = keybinds.initContainerKeys()
	ContainerKeymapBulk = keybinds.initContainerKeysBulk()
	VolumeKeymap = keybinds.initVolumeKeys()
	volumeKeymapBulk = keybinds.initVolumeKeysBulk()
}

// set tab variables, AKA IMAGES, CONTAINERS, VOLUMES, etc.
func setTabConstants(configOrder []string) TabOrderingMap {
	tabIndexMap := make(TabOrderingMap)
	for i, tab := range configOrder {
		tabIndexMap[tab] = tabId(i)
	}

	// we cannot let tab constants be default values (0) if they are not supplied in config, otherwise it will interfere with tab updation and navigation
	if index, ok := tabIndexMap["images"]; ok {
		IMAGES = index
	} else {
		IMAGES = 999
	}

	if index, ok := tabIndexMap["containers"]; ok {
		CONTAINERS = index
	} else {
		CONTAINERS = 999
	}

	if index, ok := tabIndexMap["volumes"]; ok {
		VOLUMES = index
	} else {
		VOLUMES = 999
	}
	return tabIndexMap
}
