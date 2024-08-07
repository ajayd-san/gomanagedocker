/*
	this file contains docker ops, they are put in a seperate file to facilitate easier testing
*/

package tui

import (
	"fmt"

	"github.com/ajayd-san/gomanagedocker/dockercmd"
)

type Operation func() error

// Calls dockercmd api to toggle start/stop container, and sends notification to `notificaitonChan`
func toggleStartStopContainer(cli dockercmd.DockerClient, containerInfo containerItem, activeTab tabId, notifcationChan chan notificationMetadata) Operation {

	return func() error {
		containerId := containerInfo.getId()
		err := cli.ToggleStartStopContainer(containerId)

		if err != nil {
			return err
		}

		// send notification
		msg := ""
		if containerInfo.getState() == "running" {
			msg = fmt.Sprintf("Stopped %s", containerId[:8])
		} else {
			msg = fmt.Sprintf("Started %s", containerId[:8])
		}

		notif := NewNotification(activeTab, listStatusMessageStyle.Render(msg))

		notifcationChan <- notif
		return nil
	}
}

// Calls dockercmd api to toggle pause/resume container, and sends notification to `notificaitonChan`
func togglePauseResumeContainer(client dockercmd.DockerClient, containerInfo containerItem, activeTab tabId, notificationChan chan notificationMetadata) Operation {
	return func() error {
		containerId := containerInfo.getId()
		err := client.TogglePauseResume(containerId)

		if err != nil {
			return err
		}

		// send notification
		msg := ""
		if containerInfo.getState() == "running" {
			msg = "Paused " + containerId[:8]
		} else {
			msg = "Resumed " + containerId[:8]
		}

		notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))

		return nil
	}
}

// Calls dockercmd api to restart container and sends notification to notificationChan
func toggleRestartContainer(client dockercmd.DockerClient, containerInfo containerItem, activeTab tabId, notificationChan chan notificationMetadata) Operation {
	return func() error {
		containerId := containerInfo.getId()
		err := client.RestartContainer(containerId)

		if err != nil {
			return err
		}

		msg := fmt.Sprintf("Restarted %s", containerId[:8])
		notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))

		return nil
	}
}
