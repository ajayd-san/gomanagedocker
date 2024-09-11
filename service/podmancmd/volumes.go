package podmancmd

import (
	it "github.com/ajayd-san/gomanagedocker/service/types"
	"github.com/containers/podman/v5/pkg/bindings/volumes"
)

func (pc *PodmanClient) ListVolumes() ([]it.VolumeSummary, error) {
	res, err := volumes.List(pc.cli, nil)

	if err != nil {
		return nil, err
	}

	return toVolumeSummaryArr(res), nil
}

func (pc *PodmanClient) PruneVolumes() (*it.VolumePruneReport, error) {
	report, err := volumes.Prune(pc.cli, nil)

	if err != nil {
		return nil, err
	}

	volumesPruned := 0

	for _, entry := range report {
		if entry.Err == nil {
			volumesPruned += 1
		}
	}

	return &it.VolumePruneReport{VolumesPruned: volumesPruned}, nil

}

func (pc *PodmanClient) DeleteVolume(id string, force bool) error {
	opts := &volumes.RemoveOptions{}
	opts = opts.WithForce(force)
	return volumes.Remove(pc.cli, id, opts)
}
