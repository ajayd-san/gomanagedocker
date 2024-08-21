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
			State:   entry.State,
			Command: strings.Join(entry.Command, " "),
			Status:  entry.Status,
		}
	}

	return res
}

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
		Status: info.State.Status,
	}

	return res
}
