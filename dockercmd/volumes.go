package dockercmd

import (
	"context"

	"github.com/docker/docker/api/types/volume"
)

func (dc DockerClient) ListVolumes() ([]*volume.Volume, error) {
	res, err := dc.cli.VolumeList(context.Background(), volume.ListOptions{})

	if err != nil {
		panic(err)
	}
	return res.Volumes, nil
}
