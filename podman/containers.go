package podman

import (
	"github.com/containers/podman/v5/libpod/define"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/domain/entities/reports"
	"github.com/containers/podman/v5/pkg/domain/entities/types"
	"github.com/containers/podman/v5/pkg/specgen"
)

func (p *PodmanClient) ContainerList(showContainerSize bool) ([]types.ListContainer, error) {
	opts := containers.ListOptions{}
	return containers.List(p.ctx, opts.WithSize(showContainerSize))
}

func (p *PodmanClient) ContainerInspect(id string, size bool) (*define.InspectContainerData, error) {
	opts := containers.InspectOptions{}
	opts.WithSize(size)
	return containers.Inspect(p.ctx, id, &opts)
}

func (p *PodmanClient) ContainerStart(id string) error {
	return containers.Start(p.ctx, id, nil)
}

func (p *PodmanClient) ContainerStop(id string) error {
	return containers.Stop(p.ctx, id, nil)
}

func (p *PodmanClient) ContainerRestart(id string) error {
	return containers.Restart(p.ctx, id, nil)
}

func (p *PodmanClient) ContainerPause(id string) error {
	return containers.Pause(p.ctx, id, nil)
}

func (p *PodmanClient) ContainerUnpause(id string) error {
	return containers.Unpause(p.ctx, id, nil)
}

func (p *PodmanClient) ContainerRemove(id string, removeOpts *containers.RemoveOptions) ([]*reports.RmReport, error) {
	return containers.Remove(p.ctx, id, removeOpts)
}

func (p *PodmanClient) ContainerPrune(id string, removeOpts *containers.RemoveOptions) ([]*reports.PruneReport, error) {
	return containers.Prune(p.ctx, nil)
}

func (p *PodmanClient) ContainerCreateWithSpec(spec *specgen.SpecGenerator, opts *containers.CreateOptions) (types.ContainerCreateResponse, error) {
	return containers.CreateWithSpec(p.ctx, spec, opts)
}
