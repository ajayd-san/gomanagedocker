package podman

import (
	"github.com/containers/podman/v5/pkg/bindings/images"
	"github.com/containers/podman/v5/pkg/domain/entities/reports"
	"github.com/containers/podman/v5/pkg/domain/entities/types"
	tf "github.com/containers/podman/v5/pkg/domain/entities/types"
)

func (p *PodmanClient) ImageList(opts *images.ListOptions) ([]*tf.ImageSummary, error) {
	return images.List(p.ctx, opts)
}

func (pc *PodmanClient) ImageRemove(image_ids []string, opts *images.RemoveOptions) (*tf.ImageRemoveReport, []error) {
	return images.Remove(pc.ctx, image_ids, opts)
}

func (pc *PodmanClient) ImagePrune(opts *images.PruneOptions) ([]*reports.PruneReport, error) {
	return images.Prune(pc.ctx, opts)
}

func (p *PodmanClient) ImageBuild(containerFiles []string, opts types.BuildOptions) (*tf.BuildReport, error) {
	return images.Build(p.ctx, containerFiles, opts)
}
