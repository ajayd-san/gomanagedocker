package tui

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/ajayd-san/gomanagedocker/tui/components"
	teadialog "github.com/ajayd-san/teaDialog"
	tea "github.com/charmbracelet/bubbletea"
)

type buildProgressModel struct {
	regex        *regexp.Regexp
	progressBar  components.ProgressBar
	progressChan chan string
	inner        *teadialog.InfoCard
	currentStep  string
}

func (m buildProgressModel) Init() tea.Cmd {
	return m.progressBar.Init()
}

func (m buildProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

loop:
	for {
		select {
		case status := <-m.progressChan:
			matches := m.regex.FindStringSubmatch(status)
			if matches != nil {
				currentStep, _ := strconv.ParseFloat(matches[1], 64)
				TotalSteps, _ := strconv.ParseFloat(matches[2], 64)
				m.currentStep = matches[3]

				/*
					HACK: we do `-1` since `currentStep` is the current ongoing step, it is not finished yet.
					when currentStep == TotalSteps, the progress bar would show 100% even when the build process is not finished
					which is not the right behaviour
				*/
				progressBarIncrement := (currentStep - 1) / TotalSteps

				bar, cmd := m.progressBar.Update(components.UpdateProgress(progressBarIncrement))
				cmds = append(cmds, cmd)
				m.progressBar = bar.(components.ProgressBar)
			}
		default:
			break loop
		}
	}

	bar, cmd := m.progressBar.Update(msg)
	m.progressBar = bar.(components.ProgressBar)
	m.inner.Message = fmt.Sprintf("%s\n\n%s", m.progressBar.View(), m.currentStep)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m buildProgressModel) View() string {
	return dialogContainerStyle.Render(m.inner.View())
}
