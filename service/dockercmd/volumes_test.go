package dockercmd

import (
	"testing"

	"github.com/ajayd-san/gomanagedocker/service/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/volume"
	"gotest.tools/v3/assert"
)

func TestListVolumes(t *testing.T) {
	want := []types.VolumeSummary{
		{
			Name: "1",
		},
		{
			Name: "2",
		},
		{
			Name: "3",
		},
		{
			Name: "4",
		},
	}
	vols := []*volume.Volume{
		{
			Name: "1",
		},
		{
			Name: "2",
		},
		{
			Name: "3",
		},
		{
			Name: "4",
		},
	}

	dclient := DockerClient{
		cli: &MockApi{
			mockVolumes:     vols,
			CommonAPIClient: nil,
		},
		containerListArgs: container.ListOptions{},
	}

	got, _ := dclient.ListVolumes()

	assert.DeepEqual(t, got, want)
}

func TestDeleteVolume(t *testing.T) {
	vols := []*volume.Volume{
		{
			Name: "1",
		},
		{
			Name: "2",
		},
		{
			Name: "3",
		},
		{
			Name: "4",
		},
	}

	dclient := DockerClient{
		cli: &MockApi{
			mockVolumes:     vols,
			CommonAPIClient: nil,
		},
		containerListArgs: container.ListOptions{},
	}

	dclient.DeleteVolume("1", false)

	want := vols[1:]
	finalVols := dclient.cli.(*MockApi).mockVolumes
	assert.DeepEqual(t, finalVols, want)
}
