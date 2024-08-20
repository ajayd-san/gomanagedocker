package podmancmd

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func (po *PodmanClient) InspectContainer(id string) (*types.ContainerJSON, error) {
	panic("not implemented") // TODO: Implement
}

func (po *PodmanClient) ListContainers(showContainerSize bool) []types.Container {
	return nil
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
