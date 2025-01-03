package podmancmd

import (
	"github.com/ajayd-san/gomanagedocker/podman"
	"github.com/ajayd-san/gomanagedocker/service/types"
	"github.com/containers/podman/v5/pkg/bindings/containers"
)

// TODO: investigate why we have two different containerListopts
type PodmanClient struct {
	cli               podman.PodmanAPI
	containerListOpts types.ContainerListOptions
	// internal
	listOptions containers.ListOptions
}

func (pc *PodmanClient) GetListOptions() types.ContainerListOptions {
	return pc.containerListOpts
}

func NewPodmanClient() (*PodmanClient, error) {
	api, err := podman.NewPodmanClient()

	if err != nil {
		return nil, err
	}

	return &PodmanClient{
		api,
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

func NewMockCli(cli *PodmanMockApi) *PodmanClient {
	return &PodmanClient{
		cli:               cli,
		containerListOpts: types.ContainerListOptions{},
		listOptions: containers.ListOptions{
			All: boolPtr(false),
		},
	}
}
