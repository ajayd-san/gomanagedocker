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

	for index, entry := range summary {
		res[index] = types.ContainerSummary{
			ID:         entry.ID,
			ImageID:    entry.ImageID,
			Created:    entry.Created,
			SizeRw:     entry.SizeRw,
			SizeRootFs: entry.SizeRootFs,
			Names:      entry.Names,
			State:      entry.State,
			Command:    entry.Command,
			Status:     entry.Status,
		}
	}

	return res
}
