package podmancmd

import (
	"context"

	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/docker/docker/api/types/container"
)

type PodmanClient struct {
	cli         context.Context
	listOptions containers.ListOptions
}

func (pc *PodmanClient) GetListOptions() *container.ListOptions {
	// return &pc.listOptions
	panic("unimplemented")
}

func NewPodmanClient() (*PodmanClient, error) {
	ctx, err := bindings.NewConnection(context.Background(), "unix:///run/user/1000/podman/podman.sock")

	if err != nil {
		return nil, err
	}

	return &PodmanClient{
		ctx,
		containers.ListOptions{},
	}, nil
}

// no-op since bindings.NewConnection already pings
func (pc *PodmanClient) Ping() error {
	return nil
}
