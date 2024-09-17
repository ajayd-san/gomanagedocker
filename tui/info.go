package tui

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	it "github.com/ajayd-san/gomanagedocker/service/types"
)

type InfoBoxer interface {
	InfoBox() string
}

func (im imageItem) InfoBox() string {
	var res strings.Builder
	id := strings.TrimPrefix(im.ID, "sha256:")
	id = trimToLength(id, moreInfoStyle.GetWidth())
	addEntry(&res, "id: ", id)
	addEntry(&res, "Name: ", im.getName())
	sizeInGb := float64(im.getSize())
	addEntry(&res, "Size: ", strconv.FormatFloat(sizeInGb, 'f', 2, 64)+"GB")
	if im.Containers != -1 {
		addEntry(&res, "Containers: ", strconv.Itoa(int(im.Containers)))
	}
	addEntry(&res, "Created: ", time.Unix(im.Created, 0).Format(time.UnixDate))
	return res.String()
}

func (containerInfo containerItem) InfoBox() string {
	var res strings.Builder

	id := trimToLength(containerInfo.ID, moreInfoStyle.GetWidth())
	addEntry(&res, "ID: ", id)
	addEntry(&res, "Name: ", containerInfo.getName())
	addEntry(&res, "Image: ", containerInfo.ImageName)
	addEntry(&res, "Created: ", time.Unix(containerInfo.Created, 0).Format(time.UnixDate))

	if containerInfo.Size != nil {
		log.Println("In infobox: ", containerInfo.Size)
		rootSizeInGb := float64(containerInfo.Size.RootFs) / float64(1e+9)
		SizeRwInGb := float64(containerInfo.Size.Rw) / float64(1e+9)

		addEntry(&res, "Root FS Size: ", strconv.FormatFloat(rootSizeInGb, 'f', 2, 64)+"GB")
		addEntry(&res, "SizeRw: ", strconv.FormatFloat(SizeRwInGb, 'f', 2, 64)+"GB")
	} else {
		addEntry(&res, "Root FS Size: ", "Calculating...")
		addEntry(&res, "SizeRw: ", "Calculating...")
	}

	addEntry(&res, "Command: ", containerInfo.Command)
	addEntry(&res, "State: ", containerInfo.State)

	// TODO: figure ports and mount points out
	if len(containerInfo.Mounts) > 0 {
		addEntry(&res, "Mounts: ", mountPointString(containerInfo.Mounts))
	}
	if len(containerInfo.Ports) > 0 {
		addEntry(&res, "Ports: ", portsString(containerInfo.Ports))
	}
	return res.String()
}

func (vi VolumeItem) InfoBox() string {
	var res strings.Builder

	addEntry(&res, "Name: ", vi.getName())
	addEntry(&res, "Created: ", vi.CreatedAt)
	addEntry(&res, "Driver: ", vi.Driver)

	mntPt := trimToLength(vi.Mountpoint, moreInfoStyle.GetWidth())
	addEntry(&res, "Mount Point: ", mntPt)

	if size := vi.getSize(); size != -1 {
		addEntry(&res, "Size: ", fmt.Sprintf("%f", size))
	} else {
		addEntry(&res, "Size: ", "Not Available")
	}

	return res.String()
}

func (pi PodItem) InfoBox() string {
	var res strings.Builder
	addEntry(&res, "Name: ", pi.Name)
	addEntry(&res, "ID:", pi.Id)
	addEntry(&res, "Status: ", pi.Status)
	addEntry(&res, "Containers: ", strconv.Itoa(len(pi.Containers)))

	return res.String()
}

// UTIL
func addEntry(res *strings.Builder, label string, val string) {
	label = infoEntryLabel.Render(label)
	entry := infoEntry.Render(label + val)
	res.WriteString(entry)
}

func mountPointString(mounts []string) string {
	var res strings.Builder

	// slices.SortStableFunc(mounts, func(a types.MountPoint, b types.MountPoint) int {
	// 	return cmp.Compare(a.Source, b.Source)
	// })

	for i, mount := range mounts {
		res.WriteString(mount)

		if i < len(mounts)-1 {
			res.WriteString(", ")
		}
	}

	return res.String()
}

// converts []types.Port to human readable string
func portsString(ports []it.Port) string {
	var res strings.Builder

	for _, port := range ports {
		var str string
		if port.HostPort == 0 {
			str = fmt.Sprintf("%d/%s, ", port.ContainerPort, port.Proto)
		} else {
			str = fmt.Sprintf("%d -> %d/%s, ", port.HostPort, port.ContainerPort, port.Proto)
		}
		res.WriteString(str)
	}

	return strings.TrimSuffix(res.String(), ", ")
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
