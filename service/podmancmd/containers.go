package podmancmd

import (
	"fmt"
	"os/exec"

	it "github.com/ajayd-san/gomanagedocker/service/types"
	"github.com/containers/podman/v5/pkg/bindings/containers"
)

func (pc *PodmanClient) InspectContainer(id string) (*it.InspectContainerData, error) {
	// TODO: refactor this, using `With` methods
	f := true
	raw, err := containers.Inspect(pc.cli, id, &containers.InspectOptions{
		Size: &f,
	})

	if err != nil {
		return nil, err
	}

	return &it.InspectContainerData{
		ContainerSummary: toContainerSummary(raw),
	}, nil

}

func (pc *PodmanClient) ListContainers(showContainerSize bool) []it.ContainerSummary {
	opts := pc.listOptions
	raw, err := containers.List(pc.cli, opts.WithSize(showContainerSize))

	if err != nil {
		panic(err)
	}

	return toContainerSummaryArr(raw)
}

func (pc *PodmanClient) ToggleContainerListAll() {
	if pc.containerListOpts.All {
		pc.containerListOpts.All = false
		pc.listOptions.All = boolPtr(false)
	} else {
		pc.containerListOpts.All = true
		pc.listOptions.All = boolPtr(true)
	}
}

func (pc *PodmanClient) ToggleStartStopContainer(id string, isRunning bool) error {
	var err error
	if isRunning {
		err = containers.Stop(pc.cli, id, nil)
	} else {
		err = containers.Start(pc.cli, id, nil)
	}

	return err
}

func (pc *PodmanClient) RestartContainer(id string) error {
	return containers.Restart(pc.cli, id, nil)
}

func (pc *PodmanClient) TogglePauseResume(id string, state string) error {
	var err error
	if state == "paused" {
		err = containers.Unpause(pc.cli, id, nil)
	} else if state == "running" {
		err = containers.Pause(pc.cli, id, nil)
	} else {
		err = fmt.Errorf("Cannot Pause/unPause a %s Process.", state)
	}

	return err
}

func (pc *PodmanClient) DeleteContainer(id string, opts it.ContainerRemoveOpts) error {
	podmanOpts := &containers.RemoveOptions{}
	podmanOpts = podmanOpts.WithIgnore(true)

	if opts.Force {
		podmanOpts = podmanOpts.WithForce(true)
	}
	if opts.RemoveVolumes {
		podmanOpts = podmanOpts.WithVolumes(true)
	}

	_, err := containers.Remove(pc.cli, id, podmanOpts)
	return err
}

func (po *PodmanClient) PruneContainers() (it.ContainerPruneReport, error) {
	report, err := containers.Prune(po.cli, nil)

	// only count successfully deleted containers
	containersDeleted := 0
	for _, entry := range report {
		if entry.Err == nil {
			containersDeleted += 1
		}
	}
	return it.ContainerPruneReport{ContainersDeleted: containersDeleted}, err
}

func (dc *PodmanClient) ExecCmd(id string) *exec.Cmd {
	return exec.Command("podman", "exec", "-it", id, "/bin/sh", "-c", "eval $(grep ^$(id -un): /etc/passwd | cut -d : -f 7-)")
}
