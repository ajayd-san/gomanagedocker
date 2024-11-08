package podmancmd

import (
	it "github.com/ajayd-san/gomanagedocker/service/types"
)

func (pc *PodmanClient) ListVolumes() ([]it.VolumeSummary, error) {
	res, err := pc.cli.VolumesList(nil)

	if err != nil {
		return nil, err
	}

	return toVolumeSummaryArr(res), nil
}

func (pc *PodmanClient) PruneVolumes() (*it.VolumePruneReport, error) {
	report, err := pc.cli.VolumesPrune(nil)

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
	return pc.cli.VolumesRemove(id, force)
}
