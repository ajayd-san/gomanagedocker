package tui

import (
	"fmt"
	"io"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func StartTUI(debug bool) error {
	if debug {
		f, _ := tea.LogToFile("gmd_debug.log", "debug")
		defer f.Close()
	} else {
		log.SetOutput(io.Discard)
	}

	tabs := []string{"Images", "Containers", "Volumes"}
	m := NewModel(tabs)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		return err
	}
	return nil
}
