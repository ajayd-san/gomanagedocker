package tui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ajayd-san/gomanagedocker/service/types"
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

func GetPortMappingFromStr(portStr string) ([]types.PortBinding, error) {
	portBindings := make([]types.PortBinding, 0, len(portStr))
	portStr = strings.Trim(portStr, " ")
	portMappingStrs := strings.Split(portStr, ",")

	for _, mappingStr := range portMappingStrs {
		mappingStr = strings.Trim(mappingStr, " ")
		if mappingStr == "" {
			continue
		}
		substr := strings.Split(mappingStr, ":")
		if len(substr) != 2 {
			return nil, errors.New(fmt.Sprintf("Port Mapping %s is invalid", mappingStr))
		}

		if containerPort, found := strings.CutSuffix(substr[1], "/udp"); found {
			portBindings = append(portBindings, types.PortBinding{HostPort: substr[0], ContainerPort: containerPort, Proto: "udp"})
		} else {
			containerPort, _ = strings.CutSuffix(containerPort, "/tcp")
			portBindings = append(portBindings, types.PortBinding{HostPort: substr[0], ContainerPort: containerPort, Proto: "tcp"})
		}
	}

	return portBindings, nil
}
