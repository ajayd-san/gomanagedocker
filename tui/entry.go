package tui

import (
	"fmt"
	"io"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

var POLLING_TIME time.Duration

func StartTUI(debug bool) error {
	if debug {
		f, _ := tea.LogToFile("gmd_debug.log", "debug")
		defer f.Close()
	} else {
		log.SetOutput(io.Discard)
	}

	loadConfig()

	tabs := []string{"Images", "Containers", "Volumes"}
	m := NewModel(tabs)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		return err
	}
	return nil
}

func loadConfig() {
	POLLING_TIME = viper.GetDuration("config.Polling-Time")
}
