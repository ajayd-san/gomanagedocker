/*
	this file contains docker ops, they are put in a seperate file to facilitate easier testing
*/

package tui

import (
	"fmt"
	"sync"
	"sync/atomic"

	"strings"

	"github.com/ajayd-san/gomanagedocker/service"
	"github.com/ajayd-san/gomanagedocker/service/podmancmd"
	"github.com/ajayd-san/gomanagedocker/service/types"
	"golang.design/x/clipboard"
)

type Operation func() error

// Hides/shows existed containers and sends notification to `notificationChan`
func toggleListAllContainers(client service.Service, activeTab tabId, notificationChan chan notificationMetadata) {
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
func toggleStartStopContainer(
	cli service.Service,
	containers []dockerRes,
	activeTab tabId,
	notifcationChan chan notificationMetadata,
	errChan chan error,
) Operation {
	return func() error {

		var wg sync.WaitGroup
		var successCounterStopped atomic.Uint32
		var successCounterStarted atomic.Uint32

		for _, dRes := range containers {

			wg.Add(1)
			go func() {
				containerInfo := dRes.(containerItem)
				containerId := containerInfo.GetId()

				stateStr := containerInfo.getState()
				var isRunning bool

				if stateStr == "running" {
					isRunning = true
				}
				err := cli.ToggleStartStopContainer(containerId, isRunning)

				if err != nil {
					errChan <- err
				} else {
					// send notification
					msg := ""
					if stateStr == "running" {
						msg = fmt.Sprintf("Stopped %s", containerId[:8])
						successCounterStopped.Add(1)
					} else {
						msg = fmt.Sprintf("Started %s", containerId[:8])
						successCounterStarted.Add(1)
					}

					notif := NewNotification(activeTab, listStatusMessageStyle.Render(msg))
					notifcationChan <- notif
				}
				wg.Done()
			}()
		}

		wg.Wait()

		startedContainers := successCounterStarted.Load()
		stoppedContainers := successCounterStopped.Load()

		if startedContainers+stoppedContainers > 1 {
			var msg string
			if stoppedContainers > 0 {
				msg = fmt.Sprintf("Stopped: %d", stoppedContainers)
			}
			if startedContainers > 0 {
				if msg == "" {
					msg = fmt.Sprintf("Started: %d", startedContainers)
				} else {
					msg = fmt.Sprintf("%s, Started: %d", msg, startedContainers)
				}
			}
			msg = fmt.Sprintf("%s containers", msg)
			notif := NewNotification(activeTab, listStatusMessageStyle.Render(msg))
			notifcationChan <- notif
		}
		return nil
	}
}

// Returns func that calls dockercmd api to toggle pause/resume container, and sends notification to `notificaitonChan`
func togglePauseResumeContainer(
	client service.Service,
	containers []dockerRes,
	activeTab tabId,
	notificationChan chan notificationMetadata,
	errChan chan error,
) Operation {
	return func() error {

		var wg sync.WaitGroup
		var successCounterPaused atomic.Uint32
		var successCounterResumed atomic.Uint32

		for _, container := range containers {

			wg.Add(1)
			go func() {
				containerInfo := container.(containerItem)
				containerId := containerInfo.GetId()
				state := containerInfo.getState()
				err := client.TogglePauseResume(containerId, state)

				if err != nil {
					errChan <- err
				} else {
					msg := ""
					if state == "running" {
						msg = "Paused " + containerId[:8]
						successCounterPaused.Add(1)
					} else {
						msg = "Resumed " + containerId[:8]
						successCounterResumed.Add(1)
					}
					notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))
				}
				wg.Done()
			}()
		}

		// send notification
		wg.Wait()

		resumedContainers := successCounterResumed.Load()
		pausedContainers := successCounterPaused.Load()

		if resumedContainers+pausedContainers > 1 {
			var msg string
			if pausedContainers > 0 {
				msg = fmt.Sprintf("Paused: %d", pausedContainers)
			}
			if resumedContainers > 0 {
				if msg == "" {
					msg = fmt.Sprintf("Resumed: %d", resumedContainers)
				} else {
					msg = fmt.Sprintf("%s, Resumed: %d", msg, resumedContainers)
				}
			}
			msg = fmt.Sprintf("%s containers", msg)
			notif := NewNotification(activeTab, listStatusMessageStyle.Render(msg))
			notificationChan <- notif
		}
		return nil
	}
}

// Returns func that calls dockercmd api to restart container and sends notification to notificationChan
func toggleRestartContainer(
	client service.Service,
	containers []dockerRes,
	activeTab tabId,
	notificationChan chan notificationMetadata,
	errChan chan error,
) Operation {
	return func() error {

		var wg sync.WaitGroup
		var successCounter atomic.Uint32

		for _, container := range containers {
			containerId := container.GetId()

			wg.Add(1)
			go func() {
				err := client.RestartContainer(containerId)

				if err != nil {
					errChan <- err
				} else {
					msg := fmt.Sprintf("Restarted %s", containerId[:8])
					notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))
					successCounter.Add(1)
				}

				wg.Done()
			}()
		}

		wg.Wait()

		msg := fmt.Sprintf("Restarted %d containers", successCounter.Load())
		notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))
		return nil
	}
}

// Returns func that calls dockercmd api to deletes container using `opts` as options and sends notification to notificationChan
func containerDelete(
	client service.Service,
	containerId string,
	opts types.ContainerRemoveOpts,
	activeTab tabId,
	notificationChan chan notificationMetadata,
) Operation {
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

// bulk deletes `containers`
func containerDeleteBulk(
	client service.Service,
	containers []dockerRes,
	opts types.ContainerRemoveOpts,
	activeTab tabId,
	notificationChan chan notificationMetadata,
	errChan chan error,
) Operation {
	return func() error {

		var wg sync.WaitGroup
		var successCounter atomic.Uint32

		for _, containerInfo := range containers {

			wg.Add(1)
			go func() {
				containerId := containerInfo.GetId()
				err := client.DeleteContainer(containerId, opts)

				if err != nil {
					errChan <- err
				} else {
					msg := fmt.Sprintf("Deleted %s", containerId[:8])
					notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))
					successCounter.Add(1)
				}

				wg.Done()
			}()
		}

		wg.Wait()

		if successCounter.Load() > 1 {
			msg := fmt.Sprintf("Deleted %d containers", successCounter.Load())
			notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))
		}
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
func runImage(
	client service.Service,
	config types.ContainerCreateConfig,
	activeTab tabId,
	notificationChan chan notificationMetadata,
) Operation {
	return func() error {
		_, err := client.RunImage(config)

		if err != nil {
			return err
		}

		imageId := strings.TrimPrefix(config.ImageId, "sha256:")
		notificationMsg := listStatusMessageStyle.Render(fmt.Sprintf("Run %s", imageId[:8]))

		notificationChan <- NewNotification(activeTab, notificationMsg)

		return nil
	}
}

// Deletes image with `opts` and sends notification to `notificationChan`
func imageDelete(client service.Service, imageId string, opts types.RemoveImageOptions, activeTab tabId, notificationChan chan notificationMetadata) Operation {
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

// Same as above but is a bulk operation
func imageDeleteBulk(
	client service.Service,
	items []dockerRes,
	opts types.RemoveImageOptions,
	activeTab tabId,
	notificationChan chan notificationMetadata,
	errorChan chan error,
) Operation {

	return func() error {

		var wg sync.WaitGroup
		var successCounter atomic.Uint32

		for _, item := range items {
			imageId := item.GetId()

			wg.Add(1)

			go func() {
				err := client.DeleteImage(imageId, opts)
				if err != nil {
					errorChan <- err
				} else {
					// send notification
					imageId = strings.TrimPrefix(imageId, "sha256:")
					msg := fmt.Sprintf("Deleted %s", imageId[:8])
					notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))

					successCounter.Add(1)
				}

				wg.Done()
			}()

		}

		wg.Wait()
		if successCounter.Load() > 1 {
			msg := fmt.Sprintf("Deleted %d images", successCounter.Load())
			notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))
		}

		return nil
	}
}

func volumeDeleteBulk(client service.Service, volumes []dockerRes, force bool, activeTab tabId, notificationChan chan notificationMetadata, errChan chan error) Operation {
	return func() error {
		var wg sync.WaitGroup
		var successCounter atomic.Uint32

		for _, volume := range volumes {
			wg.Add(1)
			go func() {
				volumeId := volume.GetId()
				err := client.DeleteVolume(volumeId, force)

				if err != nil {
					errChan <- err
				} else {
					msg := fmt.Sprintf("Deleted %s", volumeId[:8])
					notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))
					successCounter.Add(1)
				}

				wg.Done()
			}()
		}

		wg.Wait()

		if successCounter.Load() > 1 {
			msg := fmt.Sprintf("Deleted %d volumes", successCounter.Load())
			notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))
		}

		return nil
	}
}

func volumeDelete(client service.Service, volumeId string, force bool, activeTab tabId, notificationChan chan notificationMetadata) Operation {
	return func() error {
		err := client.DeleteVolume(volumeId, force)

		if err != nil {
			return err
		}

		notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render("Deleted"))
		return nil
	}
}

func imagePrune(client service.Service, activeTab tabId, notificationChan chan notificationMetadata) Operation {
	return func() error {

		report, err := client.PruneImages()

		if err != nil {
			return err
		}

		msg := fmt.Sprintf("Pruned %d images", report.ImagesDeleted)
		notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))
		return nil
	}
}

func volumePrune(client service.Service, activeTab tabId, notificationChan chan notificationMetadata) Operation {
	return func() error {
		report, err := client.PruneVolumes()
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("Pruned %d volumes", report.VolumesPruned)
		notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))
		return nil

	}
}

// Pods

func createPod(
	name string,
	client podmancmd.PodmanClient,
	activeTab tabId,
	notificationChan chan notificationMetadata,
) Operation {
	return func() error {
		res, err := client.CreatePod(name)
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("Created %s", res.Id[:8])
		notificationChan <- notificationMetadata{
			listId: activeTab,
			msg:    listStatusMessageStyle.Render(msg),
		}

		return nil
	}
}
func ToggleStartStopPods(
	client podmancmd.PodmanClient,
	selectedPods []dockerRes,
	activeTab tabId,
	notificationChan chan notificationMetadata,
	errChan chan error,
) Operation {
	return func() error {
		var successCounterStarted atomic.Uint32
		var successCounterStopped atomic.Uint32
		var wg sync.WaitGroup

		for _, dres := range selectedPods {
			wg.Add(1)
			go func() {
				pod := dres.(PodItem)
				isRunning := pod.Status == "running"
				err := client.ToggleStartStopPod(pod.Id, isRunning)

				if err != nil {
					errChan <- err
				} else {
					notifMsg := ""
					if isRunning {
						notifMsg = fmt.Sprintf("Stopped %s", pod.Id[:8])
						successCounterStopped.Add(1)
					} else {
						notifMsg = fmt.Sprintf("Started %s", pod.Id[:8])
						successCounterStarted.Add(1)
					}
					notificationChan <- notificationMetadata{
						listId: activeTab,
						msg:    listStatusMessageStyle.Render(notifMsg),
					}
				}

				wg.Done()
			}()
		}

		wg.Wait()

		startedPods := successCounterStarted.Load()
		stoppedPods := successCounterStopped.Load()

		if startedPods+stoppedPods > 1 {
			var msg string
			if stoppedPods > 0 {
				msg = fmt.Sprintf("Stopped: %d", stoppedPods)
			}
			if startedPods > 0 {
				if msg == "" {
					msg = fmt.Sprintf("Started: %d", startedPods)
				} else {
					msg = fmt.Sprintf("%s, Started: %d", msg, startedPods)
				}
			}
			msg = fmt.Sprintf("%s containers", msg)
			notif := NewNotification(activeTab, listStatusMessageStyle.Render(msg))
			notificationChan <- notif
		}
		return nil
	}
}

// Returns func that calls dockercmd api to toggle pause/resume container, and sends notification to `notificaitonChan`
func togglePauseResumePods(
	client *podmancmd.PodmanClient,
	pods []dockerRes,
	activeTab tabId,
	notificationChan chan notificationMetadata,
	errChan chan error,
) Operation {
	return func() error {

		var wg sync.WaitGroup
		var successCounterPaused atomic.Uint32
		var successCounterResumed atomic.Uint32

		for _, pod := range pods {

			wg.Add(1)
			go func() {
				podInfo := pod.(PodItem)
				podId := podInfo.GetId()
				state := strings.ToLower(podInfo.Status)
				err := client.TogglePauseResumePod(podId, state)

				if err != nil {
					errChan <- err
				} else {
					msg := ""
					if state == "running" {
						msg = "Paused " + podId[:8]
						successCounterPaused.Add(1)
					} else {
						msg = "Resumed " + podId[:8]
						successCounterResumed.Add(1)
					}
					notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))
				}
				wg.Done()
			}()
		}

		// send notification
		wg.Wait()

		resumedPods := successCounterResumed.Load()
		pausedPods := successCounterPaused.Load()

		if resumedPods+pausedPods > 1 {
			var msg string
			if pausedPods > 0 {
				msg = fmt.Sprintf("Paused: %d", pausedPods)
			}
			if resumedPods > 0 {
				if msg == "" {
					msg = fmt.Sprintf("Resumed: %d", resumedPods)
				} else {
					msg = fmt.Sprintf("%s, Resumed: %d", msg, resumedPods)
				}
			}
			msg = fmt.Sprintf("%s containers", msg)
			notif := NewNotification(activeTab, listStatusMessageStyle.Render(msg))
			notificationChan <- notif
		}
		return nil
	}
}

func restartPods(
	client *podmancmd.PodmanClient,
	pods []dockerRes,
	activeTab tabId,
	notificationChan chan notificationMetadata,
	errChan chan error,
) Operation {
	return func() error {

		var wg sync.WaitGroup
		var successCounter atomic.Uint32

		for _, pod := range pods {
			podId := pod.GetId()

			wg.Add(1)
			go func() {
				err := client.RestartPod(podId)

				if err != nil {
					errChan <- err
				} else {
					msg := fmt.Sprintf("Restarted %s", podId[:8])
					notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))
					successCounter.Add(1)
				}

				wg.Done()
			}()
		}

		wg.Wait()

		restartedPods := successCounter.Load()
		if restartedPods > 1 {
			msg := fmt.Sprintf("Restarted %d Pods", successCounter.Load())
			notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))
		}
		return nil
	}
}

func PodsDeleteBulk(
	client *podmancmd.PodmanClient,
	pods []dockerRes,
	activeTab tabId,
	notificationChan chan notificationMetadata,
	errChan chan error,
) Operation {
	return func() error {

		var wg sync.WaitGroup
		var successCounter atomic.Uint32

		for _, podInfo := range pods {

			wg.Add(1)
			go func() {
				containerId := podInfo.GetId()
				//TODO: use report to display number of pods deleted
				_, err := client.DeletePod(containerId, true)

				if err != nil {
					errChan <- err
				} else {
					msg := fmt.Sprintf("Deleted %s", containerId[:8])
					notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))
					successCounter.Add(1)
				}

				wg.Done()
			}()
		}

		wg.Wait()

		if successCounter.Load() > 1 {
			msg := fmt.Sprintf("Deleted %d pods", successCounter.Load())
			notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))
		}
		return nil
	}
}

func podsPrune(client *podmancmd.PodmanClient, activeTab tabId, notificationChan chan notificationMetadata) Operation {
	return func() error {
		pruneReport, err := client.PrunePods()

		if err != nil {
			return err
		}

		msg := fmt.Sprintf("Pruned %d pods", pruneReport.Removed)
		notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))
		return nil
	}
}

func podDelete(client *podmancmd.PodmanClient, id string, force bool, activeTab tabId, notificationChan chan notificationMetadata) Operation {
	return func() error {
		//TODO: use the report to display notificaiton?
		_, err := client.DeletePod(id, force)

		if err != nil {
			return err
		}

		msg := fmt.Sprintf("Deleted pod %s", id[:8])
		notificationChan <- NewNotification(activeTab, listStatusMessageStyle.Render(msg))
		return nil
	}
}
