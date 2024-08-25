package podmancmd

import (
	it "github.com/ajayd-san/gomanagedocker/service/types"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func (pc *PodmanClient) InspectContainer(id string) (*it.InspectContainerData, error) {
	// TODO: refactor this, using `With` methods
	f := true
	raw, err := containers.Inspect(pc.cli, id, &containers.InspectOptions{
		Size: &f,
	})

	if err != nil {
		return nil, err
	}

	return &it.InspectContainerData{
		ContainerSummary: toContainerSummary(raw),
	}, nil

}

func (pc *PodmanClient) ListContainers(showContainerSize bool) []it.ContainerSummary {
	opts := pc.listOptions
	raw, err := containers.List(pc.cli, opts.WithSize(showContainerSize))

	if err != nil {
		panic(err)
	}

	return toContainerSummaryArr(raw)
}

func (pc *PodmanClient) ToggleContainerListAll() {
	if pc.containerListOpts.All {
		pc.containerListOpts.All = false
		pc.listOptions.All = boolPtr(false)
	} else {
		pc.containerListOpts.All = true
		pc.listOptions.All = boolPtr(true)
	}
}

func (po *PodmanClient) ToggleStartStopContainer(id string, isRunning bool) error {
	var err error
	if isRunning {
		err = containers.Stop(po.cli, id, nil)
	} else {
		err = containers.Start(po.cli, id, nil)
	}

	return err
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
