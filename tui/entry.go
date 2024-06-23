package tui

import (
	"fmt"
	"io"
	"log"
	"time"

	config "github.com/ajayd-san/gomanagedocker/config"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/knadh/koanf/v2"
)

type TabOrderingMap map[string]tabId

var (
	IMAGES     tabId
	CONTAINERS tabId
	VOLUMES    tabId
)

var POLLING_TIME time.Duration
var CONFIG_TAB_ORDERING_SLICE []string

var globalConfig = koanf.New(".")

func StartTUI(debug bool) error {
	if debug {
		f, _ := tea.LogToFile("gmd_debug.log", "debug")
		defer f.Close()
	} else {
		log.SetOutput(io.Discard)
	}

	config.ReadConfig(globalConfig)
	loadConfig()

	m := NewModel()
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		return err
	}
	return nil
}

func loadConfig() {
	POLLING_TIME = globalConfig.Duration("config.Polling-Time")
	// I have no idea how I made this work this late in the dev process, need a reliable way to test this
	CONFIG_TAB_ORDERING_SLICE = globalConfig.Strings("config.Tab-Order")
	setTabConstants(CONFIG_TAB_ORDERING_SLICE)
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
