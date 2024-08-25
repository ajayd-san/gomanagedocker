package podmancmd

import (
	"strings"

	"github.com/ajayd-san/gomanagedocker/service/types"
	"github.com/containers/podman/v5/libpod/define"
	et "github.com/containers/podman/v5/pkg/domain/entities/types"
)

func toImageSummaryArr(summary []*et.ImageSummary) []types.ImageSummary {
	res := make([]types.ImageSummary, len(summary))

	for index, entry := range summary {
		res[index] = types.ImageSummary{
			ID:         entry.ID,
			Size:       entry.Size,
			RepoTags:   entry.RepoTags,
			Containers: int64(entry.Containers),
			Created:    entry.Created,
		}

	}

	return res
}

func toContainerSummaryArr(summary []et.ListContainer) []types.ContainerSummary {
	res := make([]types.ContainerSummary, len(summary))

	for index, entry := range summary {
		res[index] = types.ContainerSummary{
			ID:      entry.ID,
			ImageID: entry.ImageID,
			Created: entry.Created.Unix(),
			Names:   entry.Names,
			Command: strings.Join(entry.Command, " "),
			State:   entry.State,
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

func toContainerSummary(info *define.InspectContainerData) types.ContainerSummary {
	// jcart, _ := json.MarshalIndent(info, "", "\t")
	// log.Println(string(jcart))
	res := types.ContainerSummary{
		ID:      info.ID,
		ImageID: info.Image,
		Created: info.Created.Unix(),
		Names:   []string{info.Name},
		State:   info.State.Status,
		// Command:    strings.Join(entry.Command, " "),
	}

	return res
}

func toVolumeSummaryArr(entries []*et.VolumeListReport) []types.VolumeSummary {
	res := make([]types.VolumeSummary, len(entries))

	for index, entry := range entries {
		res[index] = types.VolumeSummary{
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
