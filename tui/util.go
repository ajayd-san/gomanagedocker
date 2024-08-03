package tui

import (
	"github.com/charmbracelet/bubbles/list"
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
