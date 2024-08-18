package components

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type UpdateProgress float64

type ProgressBar struct {
	Progress progress.Model
}

func (m ProgressBar) Init() tea.Cmd {
	return nil
}

func (m ProgressBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Progress.Width = msg.Width - padding*2 - 4
		if m.Progress.Width > maxWidth {
			m.Progress.Width = maxWidth
		}
		return m, nil

	case UpdateProgress:
		cmd := m.Progress.SetPercent(float64(msg))
		return m, cmd

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.Progress.Update(msg)
		m.Progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m ProgressBar) View() string {
	return m.Progress.View()

}

func NewProgressBar() ProgressBar {
	return ProgressBar{
		Progress: progress.New(progress.WithDefaultGradient()),
	}
}
