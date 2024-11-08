package podman

import (
	"github.com/containers/podman/v5/pkg/bindings/volumes"
	"github.com/containers/podman/v5/pkg/domain/entities/reports"
	"github.com/containers/podman/v5/pkg/domain/entities/types"
)

func (p *PodmanClient) VolumesList(opts *volumes.ListOptions) ([]*types.VolumeListReport, error) {
	return volumes.List(p.ctx, nil)
}

func (p *PodmanClient) VolumesRemove(id string, force bool) error {
	opts := &volumes.RemoveOptions{}
	opts = opts.WithForce(force)
	return volumes.Remove(p.ctx, id, opts)
}

func (p *PodmanClient) VolumesPrune(opts *volumes.PruneOptions) ([]*reports.PruneReport, error) {
	return volumes.Prune(p.ctx, opts)
}
