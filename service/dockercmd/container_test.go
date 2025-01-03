package dockercmd

import (
	"testing"

	it "github.com/ajayd-san/gomanagedocker/service/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"gotest.tools/v3/assert"
)

var dockerclient = NewDockerClient()

/*
was useful in deciding if i should query the containersize by default
spoiler: querying container size has huge performance impact something like +5000% time increase or something
*/
func BenchmarkContainerList(b *testing.B) {
	b.Run("Showing container size", func(b *testing.B) {
		for range b.N {
			dockerclient.ListContainers(false)
		}
	})
	b.Run("NOT Showing container size", func(b *testing.B) {
		for range b.N {
			dockerclient.ListContainers(true)
		}
	})
}

func TestListContainer(t *testing.T) {

	containers := []types.Container{
		{
			ID:         "1",
			SizeRw:     200,
			SizeRootFs: 400,
			State:      "running",
			Status:     "",
		},
		{
			ID:         "2",
			SizeRw:     201,
			SizeRootFs: 401,
			State:      "running",
		},
		{
			ID:         "3",
			SizeRw:     202,
			SizeRootFs: 402,
			State:      "created",
		},
		{
			ID:         "4",
			SizeRw:     203,
			SizeRootFs: 403,
			State:      "stopped",
		},
	}

	t.Run("Default (not showing all containers)", func(t *testing.T) {

		dclient := DockerClient{
			cli: &MockApi{
				mockContainers:  containers,
				CommonAPIClient: nil,
			},
			containerListArgs: container.ListOptions{},
		}

		got := dclient.ListContainers(false)

		want := []it.ContainerSummary{
			{
				ID: "1",
				Size: &it.SizeInfo{
					Rw:     -1,
					RootFs: -1,
				},
				State:  "running",
				Ports:  []it.Port{},
				Mounts: []string{},
			},
			{
				ID: "2",
				Size: &it.SizeInfo{
					Rw:     -1,
					RootFs: -1,
				},
				State:  "running",
				Ports:  []it.Port{},
				Mounts: []string{},
			},
		}

		assert.DeepEqual(t, want, got)
	})

	t.Run("Showing all containers", func(t *testing.T) {

		dclient := DockerClient{
			cli: &MockApi{
				mockContainers:  containers,
				CommonAPIClient: nil,
			},
			containerListArgs: container.ListOptions{},
		}

		dclient.ToggleContainerListAll()

		got := dclient.ListContainers(false)

		want := []it.ContainerSummary{
			{
				ID: "1",
				Size: &it.SizeInfo{
					Rw:     -1,
					RootFs: -1,
				},
				State:  "running",
				Ports:  []it.Port{},
				Mounts: []string{},
			},
			{
				ID: "2",
				Size: &it.SizeInfo{
					Rw:     -1,
					RootFs: -1,
				},
				State:  "running",
				Ports:  []it.Port{},
				Mounts: []string{},
			},
			{
				ID: "3",
				Size: &it.SizeInfo{
					Rw:     -1,
					RootFs: -1,
				},
				State:  "created",
				Ports:  []it.Port{},
				Mounts: []string{},
			},
			{
				ID: "4",
				Size: &it.SizeInfo{
					Rw:     -1,
					RootFs: -1,
				},
				State:  "stopped",
				Ports:  []it.Port{},
				Mounts: []string{},
			},
		}

		assert.DeepEqual(t, got, want)
	})

	t.Run("Also calculate sizes", func(t *testing.T) {

		dclient := DockerClient{
			cli: &MockApi{
				mockContainers:  containers,
				CommonAPIClient: nil,
			},
			containerListArgs: container.ListOptions{},
		}

		dclient.ToggleContainerListAll()

		got := dclient.ListContainers(true)
		want := []it.ContainerSummary{
			{
				ID: "1",
				Size: &it.SizeInfo{
					Rw:     200,
					RootFs: 400,
				},
				State:  "running",
				Ports:  []it.Port{},
				Mounts: []string{},
			},
			{
				ID: "2",
				Size: &it.SizeInfo{
					Rw:     201,
					RootFs: 401,
				},
				State:  "running",
				Ports:  []it.Port{},
				Mounts: []string{},
			},
			{
				ID: "3",
				Size: &it.SizeInfo{
					Rw:     202,
					RootFs: 402,
				},
				State:  "created",
				Ports:  []it.Port{},
				Mounts: []string{},
			},
			{
				ID: "4",
				Size: &it.SizeInfo{
					Rw:     203,
					RootFs: 403,
				},
				State:  "stopped",
				Ports:  []it.Port{},
				Mounts: []string{},
			},
		}

		assert.DeepEqual(t, got, want)
	})
}

func TestContainerToggleListAll(t *testing.T) {
	dclient := DockerClient{
		cli: &MockApi{
			mockContainers:  nil,
			CommonAPIClient: nil,
		},
		containerListArgs: container.ListOptions{},
	}

	assert.Assert(t, !dclient.containerListArgs.All)
	assert.Assert(t, !dclient.containerListOpts.All)
	dclient.ToggleContainerListAll()
	assert.Assert(t, dclient.containerListArgs.All)
	assert.Assert(t, dclient.containerListOpts.All)
}

func TestToggleStartStopContainer(t *testing.T) {
	containers := []types.Container{
		{
			ID:         "1",
			SizeRw:     200,
			SizeRootFs: 400,
			State:      "running",
			Status:     "",
		},
		{
			ID:         "2",
			SizeRw:     201,
			SizeRootFs: 401,
			State:      "running",
		},
	}

	dclient := DockerClient{
		cli: &MockApi{
			mockContainers:  containers,
			CommonAPIClient: nil,
		},
		containerListArgs: container.ListOptions{},
	}

	t.Run("Stopping container test", func(t *testing.T) {
		isRunning := true
		dclient.ToggleStartStopContainer("2", isRunning)

		state := dclient.cli.(*MockApi).mockContainers

		assert.Assert(t, state[1].State == "stopped")
	})

	t.Run("Start container test", func(t *testing.T) {
		isRunning := false
		dclient.ToggleStartStopContainer("2", isRunning)

		state := dclient.cli.(*MockApi).mockContainers

		assert.Assert(t, state[1].State == "running")
	})
}

func TestPauseUnpauseContainer(t *testing.T) {
	containers := []types.Container{
		{
			ID:    "1",
			State: "running",
		},
		{
			ID:    "2",
			State: "stopped",
		},
	}

	dclient := DockerClient{
		cli: &MockApi{
			mockContainers:  containers,
			CommonAPIClient: nil,
		},
		containerListArgs: container.ListOptions{},
	}

	t.Run("Pause running container", func(t *testing.T) {
		id := "1"
		err := dclient.TogglePauseResume(id, "running")
		assert.NilError(t, err)
		containers := dclient.cli.(*MockApi).mockContainers

		assert.Assert(t, containers[0].State == "paused")
	})

	t.Run("unpause running container", func(t *testing.T) {
		id := "1"

		err := dclient.TogglePauseResume(id, "paused")
		assert.NilError(t, err)

		containers := dclient.cli.(*MockApi).mockContainers

		assert.Assert(t, containers[0].State == "running")
	})

	t.Run("unpause stopped container(should throw error)", func(t *testing.T) {
		id := "2"
		err := dclient.TogglePauseResume(id, "exited")
		assert.ErrorContains(t, err, "Cannot Pause/unPause a")
	})
}

func TestDeleteContainer(t *testing.T) {
	containers := []types.Container{
		{
			ID:    "1",
			State: "running",
			Names: []string{"certified loverboy"},
		},
		{
			ID:    "2",
			State: "stopped",
			Names: []string{"certified *********"},
		},
	}

	dclient := DockerClient{
		cli: &MockApi{
			mockContainers:  containers,
			CommonAPIClient: nil,
		},
		containerListArgs: container.ListOptions{},
	}

	t.Run("Delete stopped container", func(t *testing.T) {
		id := "2"
		err := dclient.DeleteContainer(id, it.ContainerRemoveOpts{})
		assert.NilError(t, err)
	})

	t.Run("Try delete runing container(fails)", func(t *testing.T) {
		id := "1"
		err := dclient.DeleteContainer(id, it.ContainerRemoveOpts{})
		assert.ErrorContains(t, err, "container is running")
	})
}

func TestPruneContainer(t *testing.T) {
	containers := []types.Container{
		{
			ID:    "1",
			State: "stopped",
		},
		{
			ID:    "2",
			State: "running",
		},
		{
			ID:    "3",
			State: "stopped",
		},
	}

	dclient := DockerClient{
		cli: &MockApi{
			mockContainers:  containers,
			CommonAPIClient: nil,
		},
		containerListArgs: container.ListOptions{},
	}
	dclient.PruneContainers()

	want := []types.Container{
		{
			ID:    "2",
			State: "running",
		},
	}

	got := dclient.cli.(*MockApi).mockContainers

	assert.DeepEqual(t, want, got)
}
