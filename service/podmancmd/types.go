package podmancmd

import (
	"context"

	"github.com/ajayd-san/gomanagedocker/service/types"
	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
)

type PodmanClient struct {
	cli               context.Context
	containerListOpts types.ContainerListOptions
	// internal
	listOptions containers.ListOptions
}

func (pc *PodmanClient) GetListOptions() types.ContainerListOptions {
	return pc.containerListOpts
}

func NewPodmanClient() (*PodmanClient, error) {
	ctx, err := bindings.NewConnection(context.Background(), "unix:///run/user/1000/podman/podman.sock")

	if err != nil {
		return nil, err
	}

	return &PodmanClient{
		ctx,
		types.ContainerListOptions{},
		containers.ListOptions{
			All: boolPtr(false),
		},
	}, nil
}

// no-op since bindings.NewConnection already pings
func (pc *PodmanClient) Ping() error {
	return nil
}
