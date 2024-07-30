package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type LoadingModel struct {
	spinner      SpinnerModel
	Loaded       bool
	msg          string
	ProgressChan chan UpdateInfo
	Help         help.Model
}

type updateType int

const (
	UTLoaded updateType = iota
	UTInProgress
)

type UpdateInfo struct {
	Kind updateType
	Msg  string
}

func (m LoadingModel) Init() tea.Cmd {
	return m.spinner.Init()
}

func (m LoadingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	select {
	case update := <-m.ProgressChan:
		if update.Kind == UTLoaded {
			m.Loaded = true
		}
		m.msg = update.Msg
	default:
	}

	if !m.Loaded {
		spinner, cmd := m.spinner.Update(msg)
		m.spinner = spinner.(SpinnerModel)
		return m, cmd
	}

	return m, nil
}

func (m LoadingModel) View() string {
	return fmt.Sprintf("%s %s", m.spinner.View(), m.msg)
}

func NewLoadingModel() LoadingModel {
	return LoadingModel{
		spinner:      InitialModel(),
		ProgressChan: make(chan UpdateInfo),
		Help:         help.Model{},
	}
}
