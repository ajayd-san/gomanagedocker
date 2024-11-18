package podmancmd

import (
	"errors"
	"fmt"
	"slices"

	"github.com/containers/podman/v5/libpod/define"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/bindings/images"
	"github.com/containers/podman/v5/pkg/bindings/pods"
	"github.com/containers/podman/v5/pkg/bindings/volumes"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/containers/podman/v5/pkg/domain/entities/reports"
	"github.com/containers/podman/v5/pkg/domain/entities/types"
	"github.com/containers/podman/v5/pkg/specgen"
)

type PodmanMockApi struct {
	mockContainers []types.ListContainer
	mockVolumes    []*entities.VolumeListReport
	mockImages     []*types.ImageSummary
	mockPods       []*types.ListPodsReport
}

func (m *PodmanMockApi) SetMockImages(imgs []*types.ImageSummary) {
	m.mockImages = imgs
}

func (m *PodmanMockApi) SetMockContainers(conts []types.ListContainer) {
	m.mockContainers = conts
}

func (m *PodmanMockApi) SetMockPods(pods []*types.ListPodsReport) {
	m.mockPods = pods
}

func (m *PodmanMockApi) SetMockVolumes(vols []*entities.VolumeListReport) {
	m.mockVolumes = vols
}

// images
func (m *PodmanMockApi) ImageList(opts *images.ListOptions) ([]*types.ImageSummary, error) {
	return m.mockImages, nil
}

func (m *PodmanMockApi) ImageRemove(image_ids []string, opts *images.RemoveOptions) (*types.ImageRemoveReport, []error) {

	deleted := []string{}
	errs := []error{}

	for _, image := range image_ids {

		index := slices.IndexFunc(m.mockImages, func(i *types.ImageSummary) bool {
			if i.ID == image {
				return true
			}

			return false
		})

		if index == -1 {
			errs = append(errs, errors.New("No such image:"))
		}

		if index != -1 && m.mockImages[index].Containers > 0 {
			errs = append(errs, errors.New("unable to delete, image is in use by a container"))

		}

		if index != -1 {
			m.mockImages = slices.Delete(m.mockImages, index, index+1)
		}
		deleted = append(deleted, image)
	}

	if len(errs) > 0 {
		return &types.ImageRemoveReport{
			Deleted:  deleted,
			Untagged: []string{},
			ExitCode: 0,
		}, errs
	}

	return &types.ImageRemoveReport{
		Deleted:  deleted,
		Untagged: []string{},
		ExitCode: 0,
	}, nil

}

func (m *PodmanMockApi) ImagePrune(opts *images.PruneOptions) ([]*reports.PruneReport, error) {
	report := []*reports.PruneReport{}

	final := []*types.ImageSummary{}

	for _, img := range m.mockImages {
		if img.Containers == 0 {
			report = append(report, &reports.PruneReport{
				Id:   img.ID,
				Err:  nil,
				Size: uint64(img.Size),
			})
		} else {
			final = append(final, img)
		}
	}

	m.mockImages = final

	return report, nil
}

func (mo *PodmanMockApi) ImageBuild(containerFiles []string, opts types.BuildOptions) (*types.BuildReport, error) {
	panic("not implemented") // TODO: Implement
}

// containers
func (m *PodmanMockApi) ContainerList(opts *containers.ListOptions) ([]types.ListContainer, error) {
	final := []types.ListContainer{}

	for _, cont := range m.mockContainers {
		if cont.State == "running" || cont.State == "paused" || *opts.All {
			if !*opts.Size {
				cont.Size.RwSize = -1
				cont.Size.RootFsSize = -1
			}

			final = append(final, cont)
		}
	}

	return final, nil
}

func (m *PodmanMockApi) ContainerInspect(id string, size bool) (*define.InspectContainerData, error) {
	index := slices.IndexFunc(m.mockContainers, func(cont types.ListContainer) bool {
		if cont.ID == id {
			return true
		}

		return false
	})

	cont := m.mockContainers[index]

	state := define.InspectContainerState{
		Status: cont.State,
	}

	switch cont.State {
	case "running":
		state.Running = true
	case "paused":
		state.Paused = true
	case "dead":
		state.Dead = true
	}

	if index != -1 {
		return &define.InspectContainerData{
			ID:         cont.ID,
			Image:      cont.Image,
			Created:    cont.Created,
			ImageName:  cont.Names[0],
			State:      &state,
			SizeRw:     &cont.Size.RwSize,
			SizeRootFs: cont.Size.RootFsSize,
		}, nil
	}

	return nil, errors.New("container does not exist")

}

func (mo *PodmanMockApi) ContainerStart(id string) error {
	index := slices.IndexFunc(mo.mockContainers, func(cont types.ListContainer) bool {
		if cont.ID == id {
			return true
		}

		return false
	})

	mo.mockContainers[index].State = "running"

	return nil
}

func (mo *PodmanMockApi) ContainerStop(id string) error {
	index := slices.IndexFunc(mo.mockContainers, func(cont types.ListContainer) bool {
		if cont.ID == id {
			return true
		}

		return false
	})

	mo.mockContainers[index].State = "stopped"

	return nil
}

func (mo *PodmanMockApi) ContainerRestart(id string) error {
	panic("not implemented") // TODO: Implement
}

func (mo *PodmanMockApi) ContainerPause(id string) error {
	index := slices.IndexFunc(mo.mockContainers, func(cont types.ListContainer) bool {
		if cont.ID == id {
			return true
		}

		return false
	})

	mo.mockContainers[index].State = "paused"

	return nil
}

func (mo *PodmanMockApi) ContainerUnpause(id string) error {
	index := slices.IndexFunc(mo.mockContainers, func(cont types.ListContainer) bool {
		if cont.ID == id {
			return true
		}

		return false
	})

	mo.mockContainers[index].State = "running"

	return nil
}

func (mo *PodmanMockApi) ContainerRemove(id string, removeOpts *containers.RemoveOptions) ([]*reports.RmReport, error) {
	index := slices.IndexFunc(mo.mockContainers, func(cont types.ListContainer) bool {
		if cont.ID == id {
			return true
		}

		return false
	})

	if index == -1 {
		return nil, errors.New(fmt.Sprintf("No such container: %s", id))
	}

	if mo.mockContainers[index].State == "running" && !(removeOpts.Force != nil && *removeOpts.Force) {

		return nil, errors.New(fmt.Sprintf(
			"cannot remove container \"%s\": container is running: stop the container before removing or force remove",
			mo.mockContainers[index].Names[0],
		))
		//not exact error but works for now
	}

	mo.mockContainers = slices.Delete(mo.mockContainers, index, index+1)

	return []*reports.RmReport{{Id: id}}, nil
}

// Deletes stopped containers and returns nil, nil all the time (for now)
func (mo *PodmanMockApi) ContainerPrune() ([]*reports.PruneReport, error) {
	final := []types.ListContainer{}

	for _, cont := range mo.mockContainers {
		if cont.State == "stopped" {
			continue
		}

		final = append(final, cont)
	}

	mo.mockContainers = final
	return nil, nil

}

func (mo *PodmanMockApi) ContainerCreateWithSpec(spec *specgen.SpecGenerator, opts *containers.CreateOptions) (types.ContainerCreateResponse, error) {
	panic("not implemented") // TODO: Implement
}

// vols
func (mo *PodmanMockApi) VolumesList(opts *volumes.ListOptions) ([]*types.VolumeListReport, error) {
	return mo.mockVolumes, nil
}

func (m *PodmanMockApi) VolumesRemove(id string, force bool) error {
	final := []*types.VolumeListReport{}

	for _, vol := range m.mockVolumes {
		if vol.Name == id {
			continue
		}

		final = append(final, vol)
	}

	m.mockVolumes = final
	return nil

}

func (mo *PodmanMockApi) VolumesPrune(opts *volumes.PruneOptions) ([]*reports.PruneReport, error) {
	panic("not implemented") // TODO: Implement
}

// pods
func (po *PodmanMockApi) PodsCreate(name string) (*types.PodCreateReport, error) {
	panic("not implemented")
}

func (m *PodmanMockApi) PodsList(opts *pods.ListOptions) ([]*types.ListPodsReport, error) {
	return m.mockPods, nil
}

func (mo *PodmanMockApi) PodsRestart(id string, opts *pods.RestartOptions) (*types.PodRestartReport, error) {
	panic("not implemented") // TODO: Implement
}

func (m *PodmanMockApi) PodsPrune(opts *pods.PruneOptions) ([]*types.PodPruneReport, error) {
	final := []*types.ListPodsReport{}

	report := []*types.PodPruneReport{}
	for _, pod := range m.mockPods {
		if pod.Status == "exited" {
			report = append(report, &types.PodPruneReport{
				Err: nil,
				Id:  pod.Id,
			})
			continue
		}

		final = append(final, pod)
	}

	return report, nil
}

func (m *PodmanMockApi) PodsStop(id string, opts *pods.StopOptions) (*types.PodStopReport, error) {
	index := slices.IndexFunc(m.mockPods, func(pod *types.ListPodsReport) bool {
		if pod.Id == id {
			return true
		}

		return false
	})
	report := types.PodStopReport{}

	if index != -1 {
		report.Id = id
		m.mockContainers[index].State = "stopped"
	} else {
		report.Errs = append(report.Errs, errors.New("pod not found"))
	}

	return &report, nil
}

func (m *PodmanMockApi) PodsStart(id string, opts *pods.StartOptions) (*types.PodStartReport, error) {
	index := slices.IndexFunc(m.mockPods, func(pod *types.ListPodsReport) bool {
		if pod.Id == id {
			return true
		}

		return false
	})
	report := types.PodStartReport{}

	if index != -1 {
		report.Id = id
		m.mockContainers[index].State = "running"
	} else {
		report.Errs = append(report.Errs, errors.New("pod not found"))
	}

	return &report, nil
}

func (m *PodmanMockApi) PodsUnpause(id string, opts *pods.UnpauseOptions) (*types.PodUnpauseReport, error) {
	index := slices.IndexFunc(m.mockPods, func(pod *types.ListPodsReport) bool {
		if pod.Id == id {
			return true
		}

		return false
	})
	report := types.PodUnpauseReport{}

	if index != -1 {
		report.Id = id
		m.mockContainers[index].State = "running"
	} else {
		report.Errs = append(report.Errs, errors.New("pod not found"))
	}

	return &report, nil
}

func (m *PodmanMockApi) PodsPause(id string, opts *pods.PauseOptions) (*types.PodPauseReport, error) {
	index := slices.IndexFunc(m.mockPods, func(pod *types.ListPodsReport) bool {
		if pod.Id == id {
			return true
		}

		return false
	})
	report := types.PodPauseReport{}

	if index != -1 {
		report.Id = id
		m.mockContainers[index].State = "paused"
	} else {
		report.Errs = append(report.Errs, errors.New("pod not found"))
	}

	return &report, nil
}

func (m *PodmanMockApi) PodsRemove(id string, opts *pods.RemoveOptions) (*types.PodRmReport, error) {
	index := slices.IndexFunc(m.mockPods, func(pod *types.ListPodsReport) bool {
		if pod.Id == id {
			return true
		}

		return false
	})
	report := types.PodRmReport{}

	if index != -1 {
		if m.mockPods[index].Status == "stopped" || *opts.Force {
			report.Id = id
			m.mockPods = slices.Delete(m.mockPods, index, index+1)
		} else {
			report.Err = errors.New("running or paused containers cannot be removed without force")
		}
	}
	return &report, nil
}
