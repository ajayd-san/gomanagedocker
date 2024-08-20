package podmancmd

import (
	"github.com/ajayd-san/gomanagedocker/service/types"
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
