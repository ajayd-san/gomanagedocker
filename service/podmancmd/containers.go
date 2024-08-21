package podmancmd

import (
	it "github.com/ajayd-san/gomanagedocker/service/types"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func (pc *PodmanClient) InspectContainer(id string) (*it.InspectContainerData, error) {
	f := true
	raw, err := containers.Inspect(pc.cli, id, &containers.InspectOptions{
		Size: &f,
	})

	if err != nil {
		return nil, err
	}

	return &it.InspectContainerData{
		toContainerSummary(raw),
	}, nil

}

func (pc *PodmanClient) ListContainers(showContainerSize bool) []it.ContainerSummary {
	raw, err := containers.List(pc.cli, &pc.listOptions)

	if err != nil {
		panic(err)
	}

	return toContainerSummaryArr(raw)
}

func (po *PodmanClient) ToggleContainerListAll() {
	panic("not implemented") // TODO: Implement
}

func (po *PodmanClient) ToggleStartStopContainer(id string) error {
	panic("not implemented") // TODO: Implement
}

func (po *PodmanClient) RestartContainer(id string) error {
	panic("not implemented") // TODO: Implement
}

func (po *PodmanClient) TogglePauseResume(id string) error {
	panic("not implemented") // TODO: Implement
}

func (po *PodmanClient) DeleteContainer(id string, opts container.RemoveOptions) error {
	panic("not implemented") // TODO: Implement
}

func (po *PodmanClient) PruneContainers() (types.ContainersPruneReport, error) {
	panic("not implemented") // TODO: Implement
}
