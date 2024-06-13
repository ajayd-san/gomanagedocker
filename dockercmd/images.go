package dockercmd

import (
	"context"
	"os/exec"
	"regexp"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
)

func (dc *DockerClient) ListImages() []image.Summary {
	images, err := dc.cli.ImageList(context.Background(), image.ListOptions{ContainerCount: true})

	if err != nil {
		panic(err)
	}

	return images
}

func (dc *DockerClient) DeleteImage(id string, opts image.RemoveOptions) error {
	_, err := dc.cli.ImageRemove(context.Background(), id, opts)
	return err
}

func (dc *DockerClient) PruneImages() (types.ImagesPruneReport, error) {
	report, err := dc.cli.ImagesPrune(context.Background(), filters.Args{})
	return report, err
}

// runs docker scout and parses the output using regex
func (dc *DockerClient) ScoutImage(imageName string) (*ScoutData, error) {
	res, err := runDockerScout(imageName)

	if err != nil {
		return nil, err
	}

	return parseDockerScoutOutput(res), nil
}

// this parses docker scout quickview output
func parseDockerScoutOutput(reader []byte) *ScoutData {

	unifiedRegex := regexp.MustCompile(`\s*([\w ]+?)\s*│\s*([\w[:punct:]]+)\s*│\s+(\d)C\s+(\d+)H\s+(\d+)M\s+(\d+)L\s*(:?(\d+)\?)?`)

	matches := unifiedRegex.FindAllSubmatch(reader, -1)

	vulnerabilityEntries := make([]ImageVulnerabilities, 0, len(matches))

	for _, match := range matches {
		vulnerabilityEntries = append(vulnerabilityEntries, makeImageVulnerabilities(match))
	}

	return &ScoutData{
		ImageVulEntries: vulnerabilityEntries,
	}
}

func runDockerScout(imageId string) ([]byte, error) {
	cmd := exec.Command("docker", "scout", "quickview", imageId)

	output, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	return output, nil
}
