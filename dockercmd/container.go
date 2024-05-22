package dockercmd

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

func (dc *DockerClient) ListContainers() []types.Container {
	containers, err := dc.cli.ContainerList(context.Background(), dc.containerListArgs)

	if err != nil {
		panic(err)
	}

	return containers
}

// Toggles listing of inactive containers
func (dc *DockerClient) ToggleContainerListAll() {
	dc.containerListArgs.All = !dc.containerListArgs.All
}

// Toggles running state of container
func (dc *DockerClient) ToggleStartStopContainer(id string) error {
	info, err := dc.cli.ContainerInspect(context.Background(), id)
	if err != nil {
		return err
	}

	if info.State.Running {
		return dc.cli.ContainerStop(context.Background(), id, container.StopOptions{})
	} else {
		return dc.cli.ContainerStart(context.Background(), id, container.StartOptions{})
	}
}

// Deletes the container
func (dc *DockerClient) DeleteContainer(id string, opts container.RemoveOptions) error {
	return dc.cli.ContainerRemove(context.Background(), id, opts)
}

func (dc *DockerClient) PruneContainers() (types.ContainersPruneReport, error) {
	report, err := dc.cli.ContainersPrune(context.Background(), filters.Args{})

	if err != nil {
		return types.ContainersPruneReport{}, err
	}

	return report, nil
}

// gets logs
func (dc *DockerClient) GetContainerLogs(id string) (io.ReadCloser, error) {
	rc, err := dc.cli.ContainerLogs(context.Background(), id, container.LogsOptions{})

	if err != nil {
		return nil, err
	}

	return rc, nil
}
