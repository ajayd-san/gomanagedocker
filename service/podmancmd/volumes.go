package podmancmd

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/volume"
)

func (po *PodmanClient) ListVolumes() ([]*volume.Volume, error) {
	return []*volume.Volume{}, nil
}

func (po *PodmanClient) PruneVolumes() (*types.VolumesPruneReport, error) {
	panic("not implemented") // TODO: Implement
}

func (po *PodmanClient) DeleteVolume(id string, force bool) error {
	panic("not implemented") // TODO: Implement
}
