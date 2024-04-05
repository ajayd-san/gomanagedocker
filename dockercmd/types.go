package dockercmd

import (
	"context"

	"github.com/docker/docker/client"
)

type DockerClient struct {
	cli *client.Client
	ctx context.Context
}

func NewDockerClient() DockerClient {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return DockerClient{
		cli: cli,
		ctx: context.Background(),
	}
}
