package dockercmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"slices"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	dimage "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

type MockApi struct {
	mockContainers []types.Container
	mockVolumes    []*volume.Volume
	mockImages     []dimage.Summary
	client.CommonAPIClient
}

func (mo *MockApi) ContainerInspectWithRaw(ctx context.Context, container string, getSize bool) (types.ContainerJSON, []byte, error) {

	index := slices.IndexFunc(mo.mockContainers, func(cont types.Container) bool {
		if cont.ID == container {
			return true
		}

		return false
	})

	base, _ := mo.ContainerInspect(ctx, container)

	if getSize {
		base.SizeRw = &mo.mockContainers[index].SizeRw
		base.SizeRootFs = &mo.mockContainers[index].SizeRootFs
	}

	return base, nil, nil
}

func (m *MockApi) SetMockImages(imgs []dimage.Summary) {
	m.mockImages = imgs
}

func (m *MockApi) SetMockVolumes(vols []*volume.Volume) {
	m.mockVolumes = vols
}

func (m *MockApi) SetMockContainers(conts []types.Container) {
	m.mockContainers = conts
}

func (mo *MockApi) ContainerInspect(ctx context.Context, container string) (types.ContainerJSON, error) {
	index := slices.IndexFunc(mo.mockContainers, func(cont types.Container) bool {
		if cont.ID == container {
			return true
		}

		return false
	})

	cur := mo.mockContainers[index]

	state := types.ContainerState{}
	if cur.State == "running" {
		state.Running = true
	} else if cur.State == "paused" {
		state.Paused = true
	}

	return types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			ID:    cur.ID,
			State: &state,
		},
	}, nil

}

func (m *MockApi) ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
	final := []types.Container{}

	for _, cont := range m.mockContainers {
		if cont.State == "running" || cont.State == "paused" || options.All {
			if !options.Size {
				cont.SizeRw = -1
				cont.SizeRootFs = -1
			}

			final = append(final, cont)
		}
	}

	return final, nil
}

func (mo *MockApi) ContainerLogs(ctx context.Context, container string, options container.LogsOptions) (io.ReadCloser, error) {
	panic("not implemented") // TODO: Implement
}

func (mo *MockApi) ContainerPause(ctx context.Context, container string) error {

	index := slices.IndexFunc(mo.mockContainers, func(cont types.Container) bool {
		if cont.ID == container {
			return true
		}

		return false
	})

	mo.mockContainers[index].State = "paused"

	return nil
}

func (mo *MockApi) ContainerRemove(ctx context.Context, container string, options container.RemoveOptions) error {

	index := slices.IndexFunc(mo.mockContainers, func(cont types.Container) bool {
		if cont.ID == container {
			return true
		}

		return false
	})

	if index == -1 {
		return errors.New(fmt.Sprintf("No such container: %s", container))
	}

	if mo.mockContainers[index].State == "running" && !options.Force {
		//not exact error but works for now
		return errors.New(fmt.Sprintf(
			"cannot remove container \"%s\": container is running: stop the container before removing or force remove",
			mo.mockContainers[index].Names[0],
		))
	}

	mo.mockContainers = slices.Delete(mo.mockContainers, index, index+1)

	return nil
}

func (mo *MockApi) ContainerRestart(ctx context.Context, container string, options container.StopOptions) error {
	panic("not implemented") // TODO: Implement
}
func (mo *MockApi) ContainerStart(ctx context.Context, container string, options container.StartOptions) error {
	index := slices.IndexFunc(mo.mockContainers, func(cont types.Container) bool {
		if cont.ID == container {
			return true
		}

		return false
	})

	mo.mockContainers[index].State = "running"

	return nil
}

func (mo *MockApi) ContainerStop(ctx context.Context, container string, options container.StopOptions) error {
	index := slices.IndexFunc(mo.mockContainers, func(cont types.Container) bool {
		if cont.ID == container {
			return true
		}

		return false
	})

	mo.mockContainers[index].State = "stopped"

	return nil
}

func (mo *MockApi) ContainerUnpause(ctx context.Context, container string) error {

	index := slices.IndexFunc(mo.mockContainers, func(cont types.Container) bool {
		if cont.ID == container {
			return true
		}

		return false
	})

	mo.mockContainers[index].State = "running"

	return nil
}

func (mo *MockApi) ContainersPrune(ctx context.Context, pruneFilters filters.Args) (types.ContainersPruneReport, error) {

	final := []types.Container{}

	for _, cont := range mo.mockContainers {
		if cont.State == "stopped" {
			continue
		}

		final = append(final, cont)
	}

	mo.mockContainers = final
	return types.ContainersPruneReport{}, nil

}

func (m *MockApi) VolumeList(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error) {
	return volume.ListResponse{
		Volumes: m.mockVolumes,
	}, nil
}

func (m *MockApi) VolumeRemove(ctx context.Context, volumeID string, force bool) error {
	final := []*volume.Volume{}

	for _, vol := range m.mockVolumes {
		if vol.Name == volumeID {
			continue
		}

		final = append(final, vol)
	}

	m.mockVolumes = final
	return nil

}

func (mo *MockApi) VolumesPrune(ctx context.Context, pruneFilter filters.Args) (types.VolumesPruneReport, error) {
	panic("not implemented") // TODO: Implement
}

func (mo *MockApi) ImageBuild(ctx context.Context, context io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	newImg := dimage.Summary{
		ID:       randStr(10),
		RepoTags: options.Tags,
	}

	mo.mockImages = append(mo.mockImages, newImg)
	return types.ImageBuildResponse{
		Body:   io.NopCloser(strings.NewReader("built image!")),
		OSType: "linux",
	}, nil
}

func (m *MockApi) ImageList(ctx context.Context, options dimage.ListOptions) ([]dimage.Summary, error) {
	return m.mockImages, nil
}

func (m *MockApi) ImageRemove(ctx context.Context, image string, options dimage.RemoveOptions) ([]dimage.DeleteResponse, error) {

	res := []dimage.DeleteResponse{}

	index := slices.IndexFunc(m.mockImages, func(i dimage.Summary) bool {
		if i.ID == image {
			return true
		}

		return false
	})

	if index == -1 {
		return nil, errors.New("No such image:")
	}

	if !options.Force && m.mockImages[index].Containers > 0 {
		return nil, errors.New(fmt.Sprintf("unable to delete %s (must be forced) - image is ...", m.mockImages[index].ID))
	}

	m.mockImages = slices.Delete(m.mockImages, index, index+1)

	return res, nil

}

func (te *MockApi) ImagesPrune(ctx context.Context, pruneFilter filters.Args) (types.ImagesPruneReport, error) {
	final := []dimage.Summary{}

	for _, img := range te.mockImages {
		if img.Containers == 0 {
			continue
		}

		final = append(final, img)
	}

	te.mockImages = final

	return types.ImagesPruneReport{}, nil

}

// util
func randStr(length uint) string {
	bytes := make([]byte, int(length))
	for i := uint(0); i < length; i++ {
		bytes[i] = byte('!' + rand.Intn('~'-'!'))
	}
	return string(bytes)
}
