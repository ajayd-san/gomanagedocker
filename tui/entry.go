package tui

import (
	"fmt"
	"io"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

type TabOrderingMap map[string]tabId

var (
	IMAGES     tabId
	CONTAINERS tabId
	VOLUMES    tabId
)

var POLLING_TIME time.Duration
var CONFIG_TAB_ORDERING_SLICE []string

func StartTUI(debug bool) error {
	if debug {
		f, _ := tea.LogToFile("gmd_debug.log", "debug")
		defer f.Close()
	} else {
		log.SetOutput(io.Discard)
	}

	loadConfig()

	m := NewModel()
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		return err
	}
	return nil
}

func loadConfig() {
	POLLING_TIME = viper.GetDuration("config.Polling-Time")
	// I have no idea how I made this work this late in the dev process, need a reliable way to test this
	CONFIG_TAB_ORDERING_SLICE = viper.GetStringSlice("config.Tab-Order")
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
