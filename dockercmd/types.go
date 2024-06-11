package dockercmd

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerClient struct {
	cli               *client.Client
	containerListArgs container.ListOptions
}

func NewDockerClient() DockerClient {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return DockerClient{
		cli: cli,
		containerListArgs: container.ListOptions{
			Size:   true,
			All:    false,
			Latest: false,
		},
	}
}

func (dc DockerClient) PingDocker() error {
	_, err := dc.cli.Ping(context.Background())
	return err
}
