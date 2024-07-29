package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type loadingModel struct {
	spinner      SpinnerModel
	loaded       bool
	msg          string
	progressChan chan UpdateInfo
	Help         help.Model
}

type updateType int

const (
	UTLoaded updateType = iota
	UTInProgress
)

type UpdateInfo struct {
	kind updateType
	msg  string
}

func (m loadingModel) Init() tea.Cmd {
	return m.spinner.Init()
}

func (m loadingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	select {
	case update := <-m.progressChan:
		if update.kind == UTLoaded {
			m.loaded = true
		}
		m.msg = update.msg
	default:
	}

	if !m.loaded {
		spinner, cmd := m.spinner.Update(msg)
		m.spinner = spinner.(SpinnerModel)
		return m, cmd
	}

	return m, nil
}

func (m loadingModel) View() string {
	return fmt.Sprintf("%s %s", m.spinner.View(), m.msg)
}

func NewLoadingModel() loadingModel {
	return loadingModel{
		spinner:      initialModel(),
		progressChan: make(chan UpdateInfo),
		Help:         help.Model{},
	}
}
