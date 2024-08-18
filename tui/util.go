package tui

import (
	"github.com/ajayd-san/gomanagedocker/tui/components/list"
	tea "github.com/charmbracelet/bubbletea"
)

type notificationMetadata struct {
	listId tabId
	msg    string
}

func NotifyList(list *list.Model, msg string) tea.Cmd {
	return list.NewStatusMessage(msg)
}

func NewNotification(id tabId, msg string) notificationMetadata {
	return notificationMetadata{
		id, msg,
	}
}
