package main

import (
	"flag"
	"fmt"
	"io"
	"log"
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
	} else {
		log.SetOutput(io.Discard)
	}

	tabs := []string{"Images", "Containers", "Volumes"}
	m := tui.NewModel(tabs)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
