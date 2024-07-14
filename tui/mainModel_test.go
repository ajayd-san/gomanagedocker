package tui

import (
	"strings"
	"sync"
	"testing"

	"github.com/ajayd-san/gomanagedocker/dockercmd"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gotest.tools/v3/assert"
)

func TestNewModel(t *testing.T) {
	CONFIG_TAB_ORDERING = []string{"images", "volumes"}

	model := NewModel()

	assert.DeepEqual(t, model.Tabs, CONFIG_TAB_ORDERING)
	assert.Equal(t, model.activeTab, tabId(0))
}

func TestFetchNewData(t *testing.T) {
	api := dockercmd.MockApi{}

	containers := []types.Container{
		{
			Names:      []string{"a"},
			ID:         "1",
			SizeRw:     1e+9,
			SizeRootFs: 2e+9,
			State:      "running",
			Status:     "",
		},
		{
			Names:      []string{"b"},
			ID:         "2",
			SizeRw:     201,
			SizeRootFs: 401,
			State:      "running",
		},
		{
			Names:      []string{"c"},
			ID:         "3",
			SizeRw:     202,
			SizeRootFs: 402,
			State:      "running",
		},
		{

			Names:      []string{"d"},
			ID:         "4",
			SizeRw:     203,
			SizeRootFs: 403,
			State:      "running",
		},
	}

	imgs := []image.Summary{
		{
			Containers: 0,
			ID:         "0",
			RepoTags:   []string{"a"},
		},

		{
			Containers: 0,
			ID:         "1",
			RepoTags:   []string{"b"},
		},
		{
			Containers: 3,
			ID:         "2",
			RepoTags:   []string{"c"},
		},
		{
			Containers: 0,
			ID:         "3",
			RepoTags:   []string{"d"},
		},
	}

	api.SetMockContainers(containers)
	api.SetMockImages(imgs)

	mockcli := dockercmd.NewMockCli(&api)

	CONTAINERS = 0
	IMAGES = 1
	VOLUMES = 2
	model := MainModel{
		dockerClient: mockcli,
		activeTab:    0,
		TabContent: []listModel{
			InitList(0),
		},
		containerSizeTracker: ContainerSizeManager{
			sizeMap: make(map[string]ContainerSize),
			mu:      &sync.Mutex{},
		},
		imageIdToNameMap: map[string]string{},
	}

	newlist := model.fetchNewData(0)

	t.Run("Containers", func(t *testing.T) {
		t.Run("Assert lists", func(t *testing.T) {
			want := containers

			assert.Equal(t, len(newlist), len(want))
			for i := range len(newlist) {
				assert.Equal(t, newlist[i].getId(), want[i].ID)
				assert.Equal(t, newlist[i].getName(), strings.Join(want[i].Names, ","))
			}
		})

		t.Run("Assert containerSizeMaps", func(t *testing.T) {
			want := map[string]ContainerSize{
				"1": {1e+9, 2e+9},
				"2": {201, 401},
				"3": {202, 402},
				"4": {203, 403},
			}

			assert.DeepEqual(t, model.containerSizeTracker.sizeMap, want, cmpopts.EquateComparable(ContainerSize{}))
		})
	})

	t.Run("Images", func(t *testing.T) {
		model.nextTab()
		assert.Equal(t, model.activeTab, IMAGES)
		newlist := model.fetchNewData(IMAGES)
		t.Run("Assert images", func(t *testing.T) {

			for i := range len(newlist) {
				img := newlist[i].(imageItem)
				assert.DeepEqual(t, img.Summary, imgs[i])
			}
		})

		t.Run("Assert imageIdToNameMap", func(t *testing.T) {
			want := map[string]string{
				"0": "a",
				"1": "b",
				"2": "c",
				"3": "d",
			}
			assert.DeepEqual(t, model.imageIdToNameMap, want)
		})

	})

}

func TestInfoBoxSize(t *testing.T) {
	api := dockercmd.MockApi{}

	containers := []types.Container{
		{
			Names:      []string{"a"},
			ID:         "1",
			SizeRw:     1e+9,
			SizeRootFs: 2e+9,
			State:      "running",
			Status:     "",
		},
	}

	api.SetMockContainers(containers)

	mockcli := dockercmd.NewMockCli(&api)

	CONTAINERS = 0
	model := MainModel{
		dockerClient: mockcli,
		activeTab:    0,
		TabContent: []listModel{
			InitList(0),
		},
	}

	t.Run("With (100 width, 100 height)", func(t *testing.T) {
		model.Update(tea.WindowSizeMsg{Width: 100, Height: 100})
		assert.Equal(t, moreInfoStyle.GetHeight(), 65)
		assert.Equal(t, moreInfoStyle.GetWidth(), 55)
	})

	t.Run("With (350 width, 200 height)", func(t *testing.T) {
		model.Update(tea.WindowSizeMsg{Width: 350, Height: 200})
		assert.Equal(t, moreInfoStyle.GetHeight(), 130)
		assert.Equal(t, moreInfoStyle.GetWidth(), 192)
	})

}
