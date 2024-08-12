/*
	this file contains docker ops, they are put in a seperate file to facilitate easier testing
*/

package tui

import (
	"fmt"

	"strings"

	"github.com/ajayd-san/gomanagedocker/dockercmd"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"golang.design/x/clipboard"
)

type Operation func() error

// Hides/shows existed containers and sends notification to `notificationChan`
func toggleListAllContainers(client *dockercmd.DockerClient, activeTab tabId, notificationChan chan notificationMetadata) {
	client.ToggleContainerListAll()
	listOpts := client.GetListOptions()

	notifMsg := ""
	if listOpts.All {
		notifMsg = "List all enabled!"
	} else {
		notifMsg = "List all disabled!"
	}

	notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(notifMsg))
}

// Returns func that calls dockercmd api to toggle start/stop container, and sends notification to `notificaitonChan`
func toggleStartStopContainer(cli dockercmd.DockerClient, containerInfo containerItem, activeTab tabId, notifcationChan chan notificationMetadata) Operation {

	return func() error {
		containerId := containerInfo.GetId()
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

// Returns func that calls dockercmd api to toggle pause/resume container, and sends notification to `notificaitonChan`
func togglePauseResumeContainer(client dockercmd.DockerClient, containerInfo containerItem, activeTab tabId, notificationChan chan notificationMetadata) Operation {
	return func() error {
		containerId := containerInfo.GetId()
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

// Returns func that calls dockercmd api to restart container and sends notification to notificationChan
func toggleRestartContainer(client dockercmd.DockerClient, containerInfo containerItem, activeTab tabId, notificationChan chan notificationMetadata) Operation {
	return func() error {
		containerId := containerInfo.GetId()
		err := client.RestartContainer(containerId)

		if err != nil {
			return err
		}

		msg := fmt.Sprintf("Restarted %s", containerId[:8])
		notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))

		return nil
	}
}

// Returns func that calls dockercmd api to deletes container using `opts` as options and sends notification to notificationChan
func containerDelete(client dockercmd.DockerClient, containerId string, opts container.RemoveOptions, activeTab tabId, notificationChan chan notificationMetadata) Operation {
	return func() error {
		err := client.DeleteContainer(containerId, opts)

		if err != nil {
			return err
		}

		msg := fmt.Sprintf("Deleted %s", containerId[:8])
		notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))
		return nil
	}
}

// Copies ID of an object to clipboard and send notification to `notificationChan`
func copyIdToClipboard(object dockerRes, activeTab tabId, notificationChan chan notificationMetadata) Operation {
	return func() error {
		id := object.GetId()
		id = strings.TrimPrefix(id, "sha256:")
		id = id[:min(len(id), 20)]
		clipboard.Write(clipboard.FmtText, []byte(id))

		notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render("ID copied!"))

		return nil
	}
}

// Runs image and sends notification to `notificationChan`
func runImage(client dockercmd.DockerClient, containerConfig *container.Config, hostConfig *container.HostConfig, containerName string, activeTab tabId, notificationChan chan notificationMetadata) Operation {
	return func() error {
		_, err := client.RunImage(containerConfig, hostConfig, containerName)

		if err != nil {
			return err
		}

		imageId := strings.TrimPrefix(containerConfig.Image, "sha256:")
		notificationMsg := listStatusMessageStyle.Render(fmt.Sprintf("Run %s", imageId[:8]))

		notificationChan <- NewNotification(activeTab, notificationMsg)

		return nil
	}
}

// Deletes image with `opts` and sends notification to `notificationChan`
func imageDelete(client dockercmd.DockerClient, imageId string, opts image.RemoveOptions, activeTab tabId, notificationChan chan notificationMetadata) Operation {
	return func() error {
		err := client.DeleteImage(imageId, opts)

		if err != nil {
			return err
		}

		// send notification
		imageId = strings.TrimPrefix(imageId, "sha256:")
		msg := fmt.Sprintf("Deleted %s", imageId[:8])

		notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))

		return nil
	}
}

func volumeDelete(client dockercmd.DockerClient, volumeId string, force bool, activeTab tabId, notificationChan chan notificationMetadata) Operation {
	return func() error {
		err := client.DeleteVolume(volumeId, force)

		if err != nil {
			return err
		}

		notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render("Deleted"))
		return nil
	}
}

func imagePrune(client dockercmd.DockerClient, activeTab tabId, notificationChan chan notificationMetadata) Operation {
	return func() error {

		report, err := client.PruneImages()

		if err != nil {
			return err
		}

		msg := fmt.Sprintf("Pruned %d images", len(report.ImagesDeleted))
		notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))
		return nil
	}
}

func volumePrune(client dockercmd.DockerClient, activeTab tabId, notificationChan chan notificationMetadata) Operation {
	return func() error {
		report, err := client.PruneVolumes()
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("Pruned %d volumes", len(report.VolumesDeleted))
		notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))
		return nil

	}
}
