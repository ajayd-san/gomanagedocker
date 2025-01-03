package dockercmd

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	"github.com/ajayd-san/gomanagedocker/service/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

func (dc *DockerClient) InspectContainer(id string) (*types.InspectContainerData, error) {
	raw, _, err := dc.cli.ContainerInspectWithRaw(context.Background(), id, true)

	if err != nil {
		return nil, err
	}

	return toContainerInspectData(&raw), nil
}

func (dc *DockerClient) ListContainers(showContainerSize bool) []types.ContainerSummary {
	listArgs := dc.containerListArgs
	listArgs.Size = showContainerSize

	containers, err := dc.cli.ContainerList(context.Background(), listArgs)

	if err != nil {
		panic(err)
	}

	return toContainerSummaryArr(containers)
}

// Toggles listing of inactive containers
func (dc *DockerClient) ToggleContainerListAll() {
	dc.containerListOpts.All = !dc.containerListOpts.All
	dc.containerListArgs.All = !dc.containerListArgs.All
}

// Toggles running state of container
func (dc *DockerClient) ToggleStartStopContainer(id string, isRunning bool) error {
	if isRunning {
		return dc.cli.ContainerStop(context.Background(), id, container.StopOptions{})
	} else {
		return dc.cli.ContainerStart(context.Background(), id, container.StartOptions{})
	}
}

func (dc *DockerClient) RestartContainer(id string) error {
	return dc.cli.ContainerRestart(context.Background(), id, container.StopOptions{})
}

func (dc *DockerClient) TogglePauseResume(id string, state string) error {
	if state == "paused" {
		err := dc.cli.ContainerUnpause(context.Background(), id)

		if err != nil {
			return err
		}
	} else if state == "running" {
		err := dc.cli.ContainerPause(context.Background(), id)
		if err != nil {
			return err
		}
	} else {
		return errors.New(fmt.Sprintf("Cannot Pause/unPause a %s Process.", state))
	}

	return nil
}

// Deletes the container
func (dc *DockerClient) DeleteContainer(id string, opts types.ContainerRemoveOpts) error {
	dockerOpts := container.RemoveOptions{}
	if opts.Force {
		dockerOpts.Force = true
	}
	if opts.RemoveVolumes {
		dockerOpts.RemoveVolumes = true
	}
	return dc.cli.ContainerRemove(context.Background(), id, dockerOpts)
}

func (dc *DockerClient) PruneContainers() (types.ContainerPruneReport, error) {
	report, err := dc.cli.ContainersPrune(context.Background(), filters.Args{})

	if err != nil {
		return types.ContainerPruneReport{}, err
	}

	return types.ContainerPruneReport{ContainersDeleted: len(report.ContainersDeleted)}, nil
}

func (dc *DockerClient) ExecCmd(id string) *exec.Cmd {
	return exec.Command("docker", "exec", "-it", id, "/bin/sh", "-c", "eval $(grep ^$(id -un): /etc/passwd | cut -d : -f 7-)")
}

func (dc *DockerClient) LogsCmd(id string) *exec.Cmd {
	return exec.Command("docker", "logs", "--follow", id)
}
