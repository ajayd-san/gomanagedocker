package podmancmd

import (
	"fmt"
	"os/exec"

	it "github.com/ajayd-san/gomanagedocker/service/types"
	"github.com/containers/podman/v5/pkg/bindings/pods"
	"github.com/containers/podman/v5/pkg/domain/entities/types"
)

func (pc *PodmanClient) ListPods() ([]*types.ListPodsReport, error) {
	pods, err := pc.cli.PodsList(nil)

	if err != nil {
		return nil, err
	}

	return pods, nil
}

func (pc *PodmanClient) PausePods(id string) error {
	_, err := pc.cli.PodsPause(id, nil)
	return err
}

func (pc *PodmanClient) ResumePods(id string) error {
	_, err := pc.cli.PodsUnpause(id, nil)
	return err
}

func (pc *PodmanClient) RestartPod(id string) error {
	_, err := pc.cli.PodsRestart(id, nil)
	return err
}

func (pc *PodmanClient) PrunePods() (*it.PodsPruneReport, error) {
	reports, err := pc.cli.PodsPrune(nil)

	if err != nil {
		return nil, err
	}

	var success int
	for _, report := range reports {
		if report.Err == nil {
			success += 1
		}
	}

	return &it.PodsPruneReport{
		Removed: success,
	}, nil
}

func (pc *PodmanClient) ToggleStartStopPod(id string, isRunning bool) error {
	var err error
	if isRunning {
		_, err = pc.cli.PodsStop(id, nil)
	} else {
		_, err = pc.cli.PodsStart(id, nil)
	}
	return err
}

func (pc *PodmanClient) TogglePauseResumePod(id string, state string) error {
	var err error
	if state == "paused" {
		_, err = pc.cli.PodsUnpause(id, nil)
	} else if state == "running" {
		_, err = pc.cli.PodsPause(id, nil)
	} else {
		err = fmt.Errorf("Cannot Pause/unPause a %s Pod.", state)
	}

	return err
}

func (pc *PodmanClient) DeletePod(id string, force bool) (*it.PodsRemoveReport, error) {
	opts := pods.RemoveOptions{}

	if force {
		opts.WithForce(true)
	}

	report, err := pc.cli.PodsRemove(id, &opts)

	if err != nil {
		return nil, err
	}

	return &it.PodsRemoveReport{
		RemovedCtrs: len(report.RemovedCtrs),
	}, nil
}

func (pc *PodmanClient) LogsCmdPods(id string) *exec.Cmd {
	return exec.Command("podman", "pod", "logs", "--follow", "--color", id)
}
