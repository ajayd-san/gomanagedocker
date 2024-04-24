package dockercmd

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func (dc *DockerClient) ListContainers() []types.Container {
	containers, err := dc.cli.ContainerList(context.Background(), container.ListOptions{})

	if err != nil {
		panic(err)
	}

	return containers
}

func (dc *DockerClient) StopContainer(id string) error {
	return dc.cli.ContainerStop(context.Background(), id, container.StopOptions{})
}
