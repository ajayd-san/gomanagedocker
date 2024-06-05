package tui

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/docker/docker/api/types"
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
	sizeInGb := float64(containerInfo.SizeRootFs) / float64(1e+9)
	addEntry(&res, "Root FS Size: ", strconv.FormatFloat(sizeInGb, 'f', 2, 64))
	addEntry(&res, "Command: ", containerInfo.Command)
	addEntry(&res, "State: ", containerInfo.State)
	addEntry(&res, "Status: ", containerInfo.Status)

	if len(containerInfo.Mounts) > 0 {
		addEntry(&res, "Mounts: ", mountPointString(containerInfo.Mounts))
	}
	return res.String()
}

// UTIL
func addEntry(res *strings.Builder, label string, val string) {
	label = infoEntryLabel.Render(label)
	entry := infoEntry.Render(label + val)
	res.WriteString(entry)
}

func mountPointString(mounts []types.MountPoint) string {

	var res strings.Builder

	slices.SortStableFunc(mounts, func(a types.MountPoint, b types.MountPoint) int {
		return cmp.Compare(a.Source, b.Source)
	})

	for i, mount := range mounts {
		res.WriteString(mount.Source)

		if i < len(mounts)-1 {
			res.WriteString(", ")
		}
	}

	return res.String()
}

func mapToString(m map[string]string) string {
	var res strings.Builder

	for key, value := range m {
		res.WriteString(fmt.Sprintf("%s: %s", key, value))
	}
	return res.String()
}
