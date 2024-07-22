package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type WindowTooSmallModel struct {
	height int
	width  int
}

func (m WindowTooSmallModel) Init() tea.Cmd {
	return nil
}

func (m WindowTooSmallModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}
	return m, nil
}

func (m WindowTooSmallModel) View() string {
	return windowTooSmallStyle.Render(fmt.Sprintf(
		"Window size too small (%d x %d)\n\n"+
			"Minimum dimensions needed - Width: 65, Height: 25\n\n"+
			"Consider going fullscreen for optimal experience.",
		m.width, m.height,
	))
}

func MakeNewWindowTooSmallModel() WindowTooSmallModel {
	return WindowTooSmallModel{}
}
