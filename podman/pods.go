package podman

import (
	"github.com/containers/podman/v5/pkg/bindings/pods"
	"github.com/containers/podman/v5/pkg/domain/entities/types"
	"github.com/containers/podman/v5/pkg/specgen"
)

func (p *PodmanClient) PodsList(opts *pods.ListOptions) ([]*types.ListPodsReport, error) {
	return pods.List(p.ctx, opts)
}

func (p *PodmanClient) PodsCreate(name string) (*types.PodCreateReport, error) {
	specGen := specgen.NewPodSpecGenerator()
	specGen.Name = name

	spec := types.PodSpec{
		PodSpecGen: *specGen,
	}

	return pods.CreatePodFromSpec(p.ctx, &spec)
}

func (p *PodmanClient) PodsRestart(id string, opts *pods.RestartOptions) (*types.PodRestartReport, error) {
	return pods.Restart(p.ctx, id, opts)
}

func (p *PodmanClient) PodsPrune(opts *pods.PruneOptions) ([]*types.PodPruneReport, error) {
	return pods.Prune(p.ctx, opts)
}

func (p *PodmanClient) PodsStop(id string, opts *pods.StopOptions) (*types.PodStopReport, error) {
	return pods.Stop(p.ctx, id, opts)
}

func (p *PodmanClient) PodsStart(id string, opts *pods.StartOptions) (*types.PodStartReport, error) {
	return pods.Start(p.ctx, id, opts)
}

func (p *PodmanClient) PodsUnpause(id string, opts *pods.UnpauseOptions) (*types.PodUnpauseReport, error) {
	return pods.Unpause(p.ctx, id, opts)
}

func (p *PodmanClient) PodsPause(id string, opts *pods.PauseOptions) (*types.PodPauseReport, error) {
	return pods.Pause(p.ctx, id, opts)
}

func (p *PodmanClient) PodsRemove(id string, opts *pods.RemoveOptions) (*types.PodRmReport, error) {
	return pods.Remove(p.ctx, id, opts)
}
