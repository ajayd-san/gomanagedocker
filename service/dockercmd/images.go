package dockercmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/pkg/archive"
)

// builds a docker image from `options` and `buildContext`
func (dc *DockerClient) BuildImage(buildContext string, options types.ImageBuildOptions) (*types.ImageBuildResponse, error) {
	dockerignoreFile, err := os.Open(filepath.Join(buildContext, ".dockerignore"))

	opts := archive.TarOptions{}
	if err == nil {
		opts.ExcludePatterns = getDockerIgnorePatterns(dockerignoreFile)
	}

	tar, err := archive.TarWithOptions(buildContext, &opts)

	if err != nil {
		return nil, err
	}
	defer tar.Close()

	res, err := dc.cli.ImageBuild(context.Background(), tar, options)

	return &res, err
}

func (dc *DockerClient) ListImages() []image.Summary {
	images, err := dc.cli.ImageList(context.Background(), image.ListOptions{ContainerCount: true})

	if err != nil {
		panic(err)
	}

	return images
}

// Runs the image and returns the container ID
func (dc *DockerClient) RunImage(containerConfig *container.Config, hostConfig *container.HostConfig, containerName string) (*string, error) {
	res, err := dc.cli.ContainerCreate(
		context.Background(),
		containerConfig,
		hostConfig,
		nil,
		nil,
		containerName,
	)

	if err != nil {
		return nil, err
	}

	err = dc.cli.ContainerStart(context.Background(), res.ID, container.StartOptions{})

	if err != nil {
		return nil, err
	}

	return &res.ID, nil
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
func (dc *DockerClient) ScoutImage(ctx context.Context, imageName string) (*ScoutData, error) {
	res, err := runDockerScout(ctx, imageName)

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

func runDockerScout(ctx context.Context, imageId string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "docker", "scout", "quickview", imageId)

	output, err := cmd.Output()

	// we the error is due to Cancel() being invoked, ignore that error
	if err != nil && !errors.Is(ctx.Err(), context.Canceled) {
		return nil, err
	}

	return output, nil
}

type PortBinding struct {
	HostPort      string
	ContainerPort string
	Proto         string
}

// UTIL
func GetPortMappingFromStr(portStr string) ([]PortBinding, error) {
	portBindings := make([]PortBinding, 0, len(portStr))
	portStr = strings.Trim(portStr, " ")
	portMappingStrs := strings.Split(portStr, ",")

	for _, mappingStr := range portMappingStrs {
		mappingStr = strings.Trim(mappingStr, " ")
		if mappingStr == "" {
			continue
		}
		substr := strings.Split(mappingStr, ":")
		if len(substr) != 2 {
			return nil, errors.New(fmt.Sprintf("Port Mapping %s is invalid", mappingStr))
		}

		if containerPort, found := strings.CutSuffix(substr[1], "/udp"); found {
			portBindings = append(portBindings, PortBinding{substr[0], containerPort, "udp"})
		} else {
			containerPort, _ = strings.CutSuffix(containerPort, "/tcp")
			portBindings = append(portBindings, PortBinding{substr[0], containerPort, "tcp"})
		}
	}

	return portBindings, nil
}
