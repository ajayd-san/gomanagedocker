package podmancmd

import (
	"fmt"
	"log"
	"strings"

	it "github.com/ajayd-san/gomanagedocker/service/types"
	"github.com/containers/common/libnetwork/types"
	"github.com/containers/podman/v5/libpod/define"
	et "github.com/containers/podman/v5/pkg/domain/entities/types"
)

func toImageSummaryArr(summary []*et.ImageSummary) []it.ImageSummary {
	res := make([]it.ImageSummary, len(summary))

	for index, entry := range summary {
		res[index] = it.ImageSummary{
			ID:         entry.ID,
			Size:       entry.Size,
			RepoTags:   entry.RepoTags,
			Containers: int64(entry.Containers),
			Created:    entry.Created,
		}

	}

	return res
}

func toContainerSummaryArr(summary []et.ListContainer) []it.ContainerSummary {
	res := make([]it.ContainerSummary, len(summary))
	log.Printf("%#v", summary)

	for index, entry := range summary {
		res[index] = it.ContainerSummary{
			ID:      entry.ID,
			ImageID: entry.ImageID,
			Created: entry.Created.Unix(),
			Names:   entry.Names,
			State:   entry.State,
			Command: strings.Join(entry.Command, " "),
			Mounts:  entry.Mounts,
			Ports:   toPort(entry.Ports),
		}

		if entry.Size != nil {
			res[index].Size = &it.SizeInfo{
				Rw:     entry.Size.RwSize,
				RootFs: entry.Size.RootFsSize,
			}
		}
	}

	return res
}

func toPort(ports []types.PortMapping) []it.Port {
	res := make([]it.Port, len(ports))

	for i, port := range ports {
		res[i] = it.Port{
			HostIP:        port.HostIP,
			HostPort:      port.HostPort,
			ContainerPort: port.ContainerPort,
			Proto:         port.Protocol,
		}
	}

	return res
}

// func mapState(state *define.InspectContainerState) *types.ContainerState {
// 	return &types.ContainerState{
// 		Status:     state.Status,
// 		Running:    state.Running,
// 		Paused:     state.Paused,
// 		Restarting: state.Restarting,
// 		OOMKilled:  state.OOMKilled,
// 		Dead:       state.Dead,
// 		Pid:        state.Pid,
// 		ExitCode:   int(state.ExitCode),
// 		Error:      state.Error,
// 	}
// }

func toContainerSummary(info *define.InspectContainerData) it.ContainerSummary {
	// jcart, _ := json.MarshalIndent(info, "", "\t")
	// log.Println(string(jcart))
	res := it.ContainerSummary{
		ID:      info.ID,
		ImageID: info.Image,
		Created: info.Created.Unix(),
		Names:   []string{info.Name},
		State:   info.State.Status,
		// Command:    strings.Join(entry.Command, " "),
	}

	return res
}

func toVolumeSummaryArr(entries []*et.VolumeListReport) []it.VolumeSummary {
	res := make([]it.VolumeSummary, len(entries))

	for index, entry := range entries {
		res[index] = it.VolumeSummary{
			Name:       entry.Name,
			CreatedAt:  entry.CreatedAt.String(),
			Driver:     entry.Driver,
			Mountpoint: entry.Mountpoint,
			UsageData:  0,
		}
	}

	return res
}

func boolPtr(b bool) *bool {
	return &b
}

func getEnvMap(envVars *[]string) (map[string]string, error) {
	res := make(map[string]string)
	for _, entry := range *envVars {
		seps := strings.Split(entry, "=")
		if len(seps) != 2 {
			return nil, fmt.Errorf("Invalid environment variable: %s", entry)
		}

		res[seps[0]] = seps[1]
	}

	return res, nil
}
