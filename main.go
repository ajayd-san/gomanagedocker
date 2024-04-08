package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ajayd-san/gomanagedocker/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	debug := flag.Bool("debug", false, "bolean value to toggle debug")
	flag.Parse()

	if *debug {
		f, _ := tea.LogToFile("debug.log", "debug")
		defer f.Close()
	}

	tabs := []string{"Images", "Containers", "Volumes"}
	tabContent := []string{"Lip Gloss Tab", "Blush Tab", "Eye Shadow Tab"}
	m := tui.Model{Tabs: tabs, TabContent: tabContent}
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
