package dockercmd

import (
	"context"

	it "github.com/ajayd-san/gomanagedocker/service/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
)

func (dc *DockerClient) ListVolumes() ([]it.VolumeSummary, error) {
	res, err := dc.cli.VolumeList(context.Background(), volume.ListOptions{})

	if err != nil {
		panic(err)
	}
	return toVolumeSummaryArr(res.Volumes), nil
}

func (dc *DockerClient) PruneVolumes() (*types.VolumesPruneReport, error) {
	res, err := dc.cli.VolumesPrune(context.Background(), filters.Args{})

	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (dc *DockerClient) DeleteVolume(id string, force bool) error {
	return dc.cli.VolumeRemove(context.Background(), id, force)
}
