package podmancmd

import (
	it "github.com/ajayd-san/gomanagedocker/service/types"
	"github.com/containers/podman/v5/pkg/bindings/images"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
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

func (po *PodmanClient) RunImage(containerConfig *container.Config, hostConfig *container.HostConfig, containerName string) (*string, error) {
	panic("not implemented") // TODO: Implement
}

func (po *PodmanClient) DeleteImage(id string, opts image.RemoveOptions) error {
	panic("not implemented") // TODO: Implement
}

func (po *PodmanClient) PruneImages() (types.ImagesPruneReport, error) {
	panic("not implemented") // TODO: Implement
}
