package dockercmd

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	it "github.com/ajayd-san/gomanagedocker/service/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
)

// builds a docker image from `options` and `buildContext`
func (dc *DockerClient) BuildImage(buildContext string, options types.ImageBuildOptions) (*it.ImageBuildReport, error) {
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

	return &it.ImageBuildReport{Body: res.Body}, err
}

func (dc *DockerClient) ListImages() []it.ImageSummary {
	images, err := dc.cli.ImageList(context.Background(), image.ListOptions{ContainerCount: true})

	if err != nil {
		panic(err)
	}

	return toImageSummaryArr(images)
}

// Runs the image and returns the container ID
func (dc *DockerClient) RunImage(config it.ContainerCreateConfig) (*string, error) {

	// this is just a list of exposed ports and is used in containerConfig
	exposedPortsContainer := make(map[nat.Port]struct{}, len(config.PortBindings))
	// this is a port mapping from host to container and is used in hostConfig
	portBindings := make(nat.PortMap)

	for _, portBind := range config.PortBindings {
		port, err := nat.NewPort(portBind.Proto, portBind.ContainerPort)
		if err != nil {
			return nil, err
		}
		exposedPortsContainer[port] = struct{}{}
		portBindings[port] = []nat.PortBinding{
			{
				HostIP:   "::1",
				HostPort: portBind.HostPort,
			},
		}
	}

	dockerConfig := container.Config{
		ExposedPorts: exposedPortsContainer,
		Env:          config.Env,
		Image:        config.ImageId,
	}

	dockerHostConfig := container.HostConfig{
		PortBindings: portBindings,
	}

	res, err := dc.cli.ContainerCreate(
		context.Background(),
		&dockerConfig,
		&dockerHostConfig,
		nil,
		nil,
		config.Name,
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

func (dc *DockerClient) DeleteImage(id string, opts it.RemoveImageOptions) error {
	dockerOpts := image.RemoveOptions{
		Force:         opts.Force,
		PruneChildren: opts.NoPrune,
	}

	_, err := dc.cli.ImageRemove(context.Background(), id, dockerOpts)
	return err
}

func (dc *DockerClient) PruneImages() (it.ImagePruneReport, error) {
	report, err := dc.cli.ImagesPrune(context.Background(), filters.Args{})

	return it.ImagePruneReport{ImagesDeleted: len(report.ImagesDeleted)}, err
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
