package dockercmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/ajayd-san/gomanagedocker/service/types"
	et "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
)

func timeBenchmark(start time.Time, msg string) {
	timeTook := time.Since(start)
	log.Println(fmt.Sprintf("%s : %s", msg, timeTook))
}

func getDockerIgnorePatterns(file io.Reader) []string {
	patterns := make([]string, 0)
	buffer := bufio.NewReader(file)

	for {
		line, err := buffer.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSuffix(line, "\n")
		patterns = append(patterns, line)
	}

	return patterns
}

func toImageSummaryArr(summary []image.Summary) []types.ImageSummary {
	res := make([]types.ImageSummary, len(summary))

	for index, entry := range summary {
		res[index] = types.ImageSummary{
			ID:         entry.ID,
			Size:       entry.Size,
			RepoTags:   entry.RepoTags,
			Containers: entry.Containers,
			Created:    entry.Created,
		}
	}

	return res
}

func toContainerSummaryArr(summary []et.Container) []types.ContainerSummary {
	res := make([]types.ContainerSummary, len(summary))

	for i, entry := range summary {
		item := types.ContainerSummary{
			ServiceKind: types.Docker,
			ID:          entry.ID,
			ImageID:     entry.ImageID,
			Created:     entry.Created,
			Names:       entry.Names,
			State:       entry.State,
			Command:     entry.Command,
			Mounts:      getMounts(entry.Mounts),
			Ports:       toPort(entry.Ports),
			//BUG: this should be set to null if entry.SizeRw are 0
			Size: &types.SizeInfo{
				Rw:     entry.SizeRw,
				RootFs: entry.SizeRootFs,
			},
		}

		res[i] = item
	}

	return res
}

func toPort(ports []et.Port) []types.Port {
	res := make([]types.Port, len(ports))

	for i, port := range ports {
		res[i] = types.Port{
			HostIP:        port.IP,
			HostPort:      port.PublicPort,
			ContainerPort: port.PrivatePort,
			Proto:         port.Type,
		}
	}

	return res
}

func getMounts(mounts []et.MountPoint) []string {

	res := make([]string, len(mounts))

	for i, mount := range mounts {
		var entry strings.Builder

		entry.WriteString(mount.Source)
		entry.WriteString(":")
		entry.WriteString(mount.Destination)

		res[i] = entry.String()
	}

	return res
}

// func mapState(state *et.ContainerState) *types.ContainerState {
// 	return &types.ContainerState{
// 		Status:     state.Status,
// 		Running:    state.Running,
// 		Paused:     state.Paused,
// 		Restarting: state.Restarting,
// 		OOMKilled:  state.OOMKilled,
// 		Dead:       state.Dead,
// 		Pid:        state.Pid,
// 		ExitCode:   state.ExitCode,
// 		Error:      state.Error,
// 	}
// }

func toContainerInspectData(info *et.ContainerJSON) *types.InspectContainerData {
	res := types.ContainerSummary{
		ServiceKind: types.Docker,
		ID:          info.ID,
		ImageID:     info.Image,
		// TODO: figure out created
		// Created:    info.Created,
		Names: []string{info.Name},
		State: info.State.Status,
		// Command:    info.Command,
	}

	if info.SizeRootFs != nil && info.SizeRw != nil {
		res.Size = &types.SizeInfo{
			Rw:     *info.SizeRw,
			RootFs: *info.SizeRootFs,
		}
	}

	return &types.InspectContainerData{ContainerSummary: res}
}

func toVolumeSummaryArr(entries []*volume.Volume) []types.VolumeSummary {
	res := make([]types.VolumeSummary, len(entries))

	for index, entry := range entries {
		res[index] = types.VolumeSummary{
			Name:       entry.Name,
			CreatedAt:  entry.CreatedAt,
			Driver:     entry.Driver,
			Mountpoint: entry.Mountpoint,
		}

		if entry.UsageData != nil {
			res[index].UsageData = entry.UsageData.Size
		}
	}

	return res
}
