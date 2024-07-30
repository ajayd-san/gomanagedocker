package tui

import (
	"github.com/ajayd-san/gomanagedocker/tui/components"
	teadialog "github.com/ajayd-san/teaDialog"
	tea "github.com/charmbracelet/bubbletea"
)

type buildProgressModel struct {
	loading components.LoadingModel
	inner   *teadialog.InfoCard
}

func (m buildProgressModel) Init() tea.Cmd {
	return m.loading.Init()
}

func (m buildProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd2 tea.Cmd
	if !m.loading.Loaded {
		spinner, cmd := m.loading.Update(msg)
		m.loading = spinner.(components.LoadingModel)
		m.inner.Message = m.loading.View()
		cmd2 = cmd
	}

	return m, cmd2
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m buildProgressModel) View() string {
	return dialogContainerStyle.Render(m.inner.View())
}
