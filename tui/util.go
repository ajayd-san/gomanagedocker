package tui

import (
	"strings"

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

func mapSpecialCharactersToString(c string) string {
	strings.ToLower(c)
	c = strings.ReplaceAll(c, "down", "↓")
	c = strings.ReplaceAll(c, "up", "↑")
	c = strings.ReplaceAll(c, "left", "←")
	c = strings.ReplaceAll(c, "right", "→")
	return c
}

func ArrayToString(arr []string) string {
	return mapSpecialCharactersToString(strings.Join(arr, "/"))
}
