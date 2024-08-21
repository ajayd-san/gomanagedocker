package podmancmd

import (
	it "github.com/ajayd-san/gomanagedocker/service/types"
	"github.com/containers/podman/v5/pkg/bindings/volumes"
	"github.com/docker/docker/api/types"
)

func (pc *PodmanClient) ListVolumes() ([]it.VolumeSummary, error) {
	res, err := volumes.List(pc.cli, nil)

	if err != nil {
		return nil, err
	}

	return toVolumeSummaryArr(res), nil
}

func (po *PodmanClient) PruneVolumes() (*types.VolumesPruneReport, error) {
	panic("not implemented") // TODO: Implement
}

func (po *PodmanClient) DeleteVolume(id string, force bool) error {
	panic("not implemented") // TODO: Implement
}
