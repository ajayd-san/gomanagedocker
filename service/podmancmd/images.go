package podmancmd

import (
	it "github.com/ajayd-san/gomanagedocker/service/types"
	"github.com/containers/podman/v5/pkg/bindings/images"
	"github.com/docker/docker/api/types"
)

func (po *PodmanClient) BuildImage(buildContext string, options types.ImageBuildOptions) (*types.ImageBuildResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (pc *PodmanClient) ListImages() []it.ImageSummary {
	raw, err := images.List(pc.cli, nil)

	if err != nil {
		panic(err)
	}

	return toImageSummaryArr(raw)
}

func (pc *PodmanClient) RunImage(config it.ContainerCreateConfig) (*string, error) {
	// spec := specgen.NewSpecGenerator(image, false)
	// containers.CreateWithSpec()
	panic("un implemented")
}

func (pc *PodmanClient) DeleteImage(id string, opts it.RemoveImageOptions) error {
	_, errs := images.Remove(pc.cli, []string{id}, &images.RemoveOptions{
		All:            &opts.All,
		Force:          &opts.Force,
		Ignore:         &opts.Ignore,
		LookupManifest: &opts.LookupManifest,
		NoPrune:        &opts.NoPrune,
	})

	if errs != nil {
		return errs[0]
	}

	return nil
}

func (pc *PodmanClient) PruneImages() (it.ImagePruneReport, error) {
	t := true
	reports, err := images.Prune(pc.cli, &images.PruneOptions{
		All: &t,
	})

	return it.ImagePruneReport{ImagesDeleted: len(reports)}, err
}
