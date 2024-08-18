package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestWindowTooSmallModel_Update(t *testing.T) {
	model := WindowTooSmallModel{}

	msg := tea.WindowSizeMsg{Width: 100, Height: 24}
	updatedModel, _ := model.Update(msg)

	if updatedModel.(WindowTooSmallModel).width != 100 {
		t.Errorf("expected width to be 100, got %d", updatedModel.(WindowTooSmallModel).width)
	}
	if updatedModel.(WindowTooSmallModel).height != 24 {
		t.Errorf("expected height to be 24, got %d", updatedModel.(WindowTooSmallModel).height)
	}
}

func TestWindowTooSmallModel_View(t *testing.T) {
	model := WindowTooSmallModel{width: 100, height: 24}
	expectedOutput := windowTooSmallStyle.Render(
		"Window size too small (100 x 24)\n\n" +
			"Minimum dimensions needed - Width: 65, Height: 25\n\n" +
			"Consider going fullscreen for optimal experience.",
	)

	if model.View() != expectedOutput {
		t.Errorf("expected %q, got %q", expectedOutput, model.View())
	}
}
