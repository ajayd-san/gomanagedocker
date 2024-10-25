/*
This package is a wrapper around podman bindings, I've done this to facilitate testing since the
podman bindings are pure functions. There is no way to swap them with stubs
*/
package podman

import (
	"context"

	"github.com/containers/podman/v5/libpod/define"
	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/bindings/images"
	"github.com/containers/podman/v5/pkg/bindings/pods"
	"github.com/containers/podman/v5/pkg/bindings/volumes"
	"github.com/containers/podman/v5/pkg/domain/entities/reports"
	"github.com/containers/podman/v5/pkg/domain/entities/types"
	tf "github.com/containers/podman/v5/pkg/domain/entities/types"
	"github.com/containers/podman/v5/pkg/specgen"
)

type PodmanAPI interface {
	//images
	ImageList(opts *images.ListOptions) ([]*tf.ImageSummary, error)
	ImageRemove(image_ids []string, opts *images.RemoveOptions) (*tf.ImageRemoveReport, []error)
	ImagePrune(opts *images.PruneOptions) ([]*reports.PruneReport, error)
	ImageBuild(containerFiles []string, opts types.BuildOptions) (*tf.BuildReport, error)

	// containers
	ContainerList(opts *containers.ListOptions) ([]types.ListContainer, error)
	ContainerInspect(id string, size bool) (*define.InspectContainerData, error)
	ContainerStart(id string) error
	ContainerStop(id string) error
	ContainerRestart(id string) error
	ContainerPause(id string) error
	ContainerUnpause(id string) error
	ContainerRemove(id string, removeOpts *containers.RemoveOptions) ([]*reports.RmReport, error)
	ContainerPrune() ([]*reports.PruneReport, error)
	ContainerCreateWithSpec(spec *specgen.SpecGenerator, opts *containers.CreateOptions) (types.ContainerCreateResponse, error)

	// vols
	VolumesList(opts *volumes.ListOptions) ([]*types.VolumeListReport, error)
	VolumesRemove(id string, force bool) error
	VolumesPrune(opts *volumes.PruneOptions) ([]*reports.PruneReport, error)

	//pods
	PodsList(opts *pods.ListOptions) ([]*types.ListPodsReport, error)
	PodsRestart(id string, opts *pods.RestartOptions) (*types.PodRestartReport, error)
	PodsPrune(opts *pods.PruneOptions) ([]*types.PodPruneReport, error)
	PodsStop(id string, opts *pods.StopOptions) (*types.PodStopReport, error)
	PodsStart(id string, opts *pods.StartOptions) (*types.PodStartReport, error)
	PodsUnpause(id string, opts *pods.UnpauseOptions) (*types.PodUnpauseReport, error)
	PodsPause(id string, opts *pods.PauseOptions) (*types.PodPauseReport, error)
	PodsRemove(id string, opts *pods.RemoveOptions) (*types.PodRmReport, error)
}

type PodmanClient struct {
	ctx context.Context
}

const defaultSocket string = "unix:///run/user/1000/podman/podman.sock"

func NewPodmanClient() (PodmanAPI, error) {
	ctx, err := bindings.NewConnection(context.Background(), defaultSocket)

	if err != nil {
		return nil, err
	}

	return &PodmanClient{
		ctx: ctx,
	}, nil
}
