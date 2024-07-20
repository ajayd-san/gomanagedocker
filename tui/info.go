package tui

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
)

func populateImageInfoBox(imageinfo imageItem) string {
	var res strings.Builder
	id := strings.TrimPrefix(imageinfo.ID, "sha256:")
	id = trimToLength(id, moreInfoStyle.GetWidth())
	addEntry(&res, "id: ", id)
	addEntry(&res, "Name: ", imageinfo.getName())
	sizeInGb := float64(imageinfo.getSize())
	addEntry(&res, "Size: ", strconv.FormatFloat(sizeInGb, 'f', 2, 64)+"GB")
	if imageinfo.Containers != -1 {
		addEntry(&res, "Containers: ", strconv.Itoa(int(imageinfo.Containers)))
	}
	addEntry(&res, "Created: ", time.Unix(imageinfo.Created, 0).Format(time.UnixDate))
	return res.String()
}

func populateVolumeInfoBox(volumeInfo VolumeItem) string {
	var res strings.Builder

	addEntry(&res, "Name: ", volumeInfo.getName())
	addEntry(&res, "Created: ", volumeInfo.CreatedAt)
	addEntry(&res, "Driver: ", volumeInfo.Driver)

	mntPt := trimToLength(volumeInfo.Mountpoint, moreInfoStyle.GetWidth())
	addEntry(&res, "Mount Point: ", mntPt)

	if size := volumeInfo.getSize(); size != -1 {
		addEntry(&res, "Size: ", fmt.Sprintf("%f", size))
	} else {
		addEntry(&res, "Size: ", "Not Available")
	}

	return res.String()
}

func populateContainerInfoBox(containerInfo containerItem, containerSizeTracker *ContainerSizeManager, imageIdToNameMap map[string]string) string {
	var res strings.Builder

	id := trimToLength(containerInfo.ID, moreInfoStyle.GetWidth())
	addEntry(&res, "ID: ", id)
	addEntry(&res, "Name: ", containerInfo.getName())
	addEntry(&res, "Image: ", imageIdToNameMap[containerInfo.ImageID])
	addEntry(&res, "Created: ", time.Unix(containerInfo.Created, 0).Format(time.UnixDate))

	//this is a pretty trivial refactor to make this look cleaner, but I'm too lazy to do this
	// whoever completes this bounty will win......nothing (except my heart)
	if mutexok := containerSizeTracker.mu.TryLock(); mutexok {
		if containerSizeInfo, ok := containerSizeTracker.sizeMap[containerInfo.ID]; ok {
			rootSizeInGb := float64(containerSizeInfo.rootFs) / float64(1e+9)
			SizeRwInGb := float64(containerSizeInfo.sizeRw) / float64(1e+9)

			addEntry(&res, "Root FS Size: ", strconv.FormatFloat(rootSizeInGb, 'f', 2, 64)+"GB")
			addEntry(&res, "SizeRw: ", strconv.FormatFloat(SizeRwInGb, 'f', 2, 64)+"GB")
		} else {
			addEntry(&res, "Root FS Size: ", "Calculating...")
			addEntry(&res, "SizeRw: ", "Calculating...")
		}

		containerSizeTracker.mu.Unlock()
	} else {
		addEntry(&res, "Root FS Size: ", "Calculating...")
		addEntry(&res, "SizeRw: ", "Calculating...")
	}

	addEntry(&res, "Command: ", containerInfo.Command)
	addEntry(&res, "State: ", containerInfo.State)

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

func trimToLength(id string, availableWidth int) string {
	return id[:min(availableWidth-10, len(id))]
}
