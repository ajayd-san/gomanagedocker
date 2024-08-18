// this is a wrapper on tea.InfoCard since I need my implementation to update
package tui

import (
	"fmt"

	teadialog "github.com/ajayd-san/teaDialog"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/ajayd-san/gomanagedocker/dockercmd"
	"github.com/ajayd-san/gomanagedocker/tui/components"
)

const customLoadingMessage = "Loading (this may take some time)..."

type DockerScoutInfoCard struct {
	tableChan  chan *TableModel
	loaded     bool
	spinner    components.SpinnerModel
	tableModel *TableModel
	inner      *teadialog.InfoCard
	f          func() (*dockercmd.ScoutData, error)
}

func (m DockerScoutInfoCard) Init() tea.Cmd {
	go func() {
		ScoutData, err := m.f()
		if err != nil {
			return
		}
		m.tableChan <- NewTable(*ScoutData)
	}()

	return m.spinner.Init()
}

func (m DockerScoutInfoCard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	if !m.loaded {
		spinner, cmd := m.spinner.Update(msg)
		m.spinner = spinner.(components.SpinnerModel)
		m.inner.Message = fmt.Sprintf("%s %s", m.spinner.View(), customLoadingMessage)
		cmds = append(cmds, cmd)

		select {
		case m.tableModel = <-m.tableChan:
			m.loaded = true
			m.inner.Message = m.tableModel.View()
		default:
		}
	}

	update, cmd := m.inner.Update(msg)
	infoCard := update.(teadialog.InfoCard)
	m.inner = &infoCard

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m DockerScoutInfoCard) View() string {
	return dialogContainerStyle.Render(m.inner.View())
}
