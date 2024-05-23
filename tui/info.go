package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
)

func PopulateInfoBox(tab tabId, item list.Item) string {
	temp, _ := item.(dockerRes)
	switch tab {
	case images:
		if it, ok := temp.(imageItem); ok {
			return populateImageInfoBox(it)
		}

	case containers:
		if ct, ok := temp.(containerItem); ok {
			return populateContainerInfoBox(ct)
		}

	case volumes:
		if vt, ok := temp.(VolumeItem); ok {
			return populateVolumeInfoBox(vt)
		}
	}
	return ""
}

func populateImageInfoBox(imageinfo imageItem) string {
	var res strings.Builder
	// res.WriteString(addEntry(res, "Containers", imageinfo.Containers))
	addEntry(&res, "id: ", strings.TrimPrefix(imageinfo.ID, "sha256:"))
	addEntry(&res, "Name: ", imageinfo.getName())
	//BUG: this always shows -1
	addEntry(&res, "Containers: ", fmt.Sprintf("%d", imageinfo.Containers))
	addEntry(&res, "Created: ", time.Unix(imageinfo.Created, 0).Format(time.UnixDate))
	return res.String()
}

func populateVolumeInfoBox(volumeInfo VolumeItem) string {
	var res strings.Builder

	addEntry(&res, "Name: ", volumeInfo.getName())
	addEntry(&res, "Created: ", volumeInfo.CreatedAt)
	addEntry(&res, "Driver: ", volumeInfo.Driver)
	addEntry(&res, "Mount Point: ", volumeInfo.Mountpoint)

	if size := volumeInfo.getSize(); size != -1 {
		addEntry(&res, "Size: ", fmt.Sprintf("%f", size))
	} else {
		addEntry(&res, "Size: ", "Not Available")
	}

	return res.String()
}

func populateContainerInfoBox(containerInfo containerItem) string {
	var res strings.Builder

	addEntry(&res, "ID: ", containerInfo.ID)
	addEntry(&res, "Name: ", containerInfo.getName())
	addEntry(&res, "Image: ", containerInfo.Image)
	addEntry(&res, "Created: ", time.Unix(containerInfo.Created, 0).Format(time.UnixDate))
	addEntry(&res, "Command: ", containerInfo.Command)
	addEntry(&res, "State: ", containerInfo.State)
	addEntry(&res, "Status: ", containerInfo.Status)
	return res.String()
}

// UTIL
func addEntry(res *strings.Builder, label string, val string) {
	label = infoEntryLabel.Render(label)
	entry := infoEntry.Render(label + val)
	res.WriteString(entry)
}

func mapToString(m map[string]string) string {
	var res strings.Builder

	for key, value := range m {
		res.WriteString(fmt.Sprintf("%s: %s", key, value))
	}
	return res.String()
}
