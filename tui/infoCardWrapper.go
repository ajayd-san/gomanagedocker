// this is a wrapper on tea.InfoCard since I need my implementation to update
package tui

import (
	"fmt"

	teadialog "github.com/ajayd-san/teaDialog"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/ajayd-san/gomanagedocker/dockercmd"
)

const customLoadingMessage = "Loading (this may take some time)..."

type InfoCardWrapperModel struct {
	tableChan  chan *TableModel
	loaded     bool
	spinner    SpinnerModel
	tableModel *TableModel
	inner      *teadialog.InfoCard
	f          func() (*dockercmd.ScoutData, error)
}

func (m InfoCardWrapperModel) Init() tea.Cmd {
	go func() {
		ScoutData, err := m.f()
		if err != nil {
			return
		}
		m.tableChan <- NewTable(*ScoutData)
	}()

	return m.spinner.Init()
}

func (m InfoCardWrapperModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd2 tea.Cmd
	if !m.loaded {
		spinner, cmd := m.spinner.Update(msg)
		m.spinner = spinner.(SpinnerModel)
		m.inner.Message = fmt.Sprintf("%s %s", m.spinner.View(), customLoadingMessage)
		cmd2 = cmd

		select {
		case m.tableModel = <-m.tableChan:
			m.loaded = true
			m.inner.Message = m.tableModel.View()
		default:
		}
	}

	return m, cmd2
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m InfoCardWrapperModel) View() string {
	return dialogContainerStyle.Render(m.inner.View())
}
