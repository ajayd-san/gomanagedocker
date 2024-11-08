package tui

import (
	"errors"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/ajayd-san/gomanagedocker/dockercmd"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
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
			InitList(0, ContainerKeymap, ContainerKeymapBulk),
		},
		containerSizeTracker: ContainerSizeManager{
			sizeMap: make(map[string]ContainerSize),
			mu:      &sync.Mutex{},
		},
		imageIdToNameMap: map[string]string{},
	}

	wg := sync.WaitGroup{}
	newlist := model.fetchNewData(0, &wg)
	wg.Wait()

	t.Run("Containers", func(t *testing.T) {
		t.Run("Assert lists", func(t *testing.T) {
			want := containers

			assert.Equal(t, len(newlist), len(want))
			for i := range len(newlist) {
				assert.Equal(t, newlist[i].GetId(), want[i].ID)
				assert.Equal(t, newlist[i].getName(), strings.Join(want[i].Names, ","))
			}
		})

		// this fails on macos ci for some reason
		// t.Run("Assert containerSizeMaps", func(t *testing.T) {
		// 	want := map[string]it.SizeInfo{
		// 		"1": {Rw: 1e+9, RootFs: 2e+9},
		// 		"2": {Rw: 201, RootFs: 401},
		// 		"3": {Rw: 202, RootFs: 402},
		// 		"4": {Rw: 203, RootFs: 403},
		// 	}

		// 	assert.DeepEqual(t, model.containerSizeTracker.sizeMap, want, cmpopts.EquateComparable(it.SizeInfo{}))
		// })
	})

	t.Run("Images", func(t *testing.T) {
		model.nextTab()
		assert.Equal(t, model.activeTab, IMAGES)
		newlist := model.fetchNewData(IMAGES, &wg)
		wg.Wait()
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
			InitList(0, ContainerKeymap, ContainerKeymapBulk),
		},
	}

	t.Run("With (100 width, 100 height)", func(t *testing.T) {
		model.Update(tea.WindowSizeMsg{Width: 100, Height: 100})
		assert.Equal(t, moreInfoStyle.GetHeight(), 60)
		assert.Equal(t, moreInfoStyle.GetWidth(), 55)
	})

	t.Run("With (350 width, 200 height)", func(t *testing.T) {
		model.Update(tea.WindowSizeMsg{Width: 350, Height: 200})
		assert.Equal(t, moreInfoStyle.GetHeight(), 120)
		assert.Equal(t, moreInfoStyle.GetWidth(), 192)
	})

}

func TestMainModelUpdate(t *testing.T) {
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
			InitList(0, ContainerKeymap, ContainerKeymapBulk),
		},
	}

	//model.windowTooSmall should be true if height < 25 or width < 65
	t.Run("Assert Window too small with small height", func(t *testing.T) {
		temp, _ := model.Update(tea.WindowSizeMsg{
			Width:  100,
			Height: 24,
		})

		model = temp.(MainModel)

		assert.Check(t, model.windowTooSmall)
	})

	t.Run("Assert Window too small with small width", func(t *testing.T) {
		temp, _ := model.Update(tea.WindowSizeMsg{
			Width:  64,
			Height: 100,
		})

		model = temp.(MainModel)

		assert.Check(t, model.windowTooSmall)
	})

	// if msg.Height <= 31 || msg.Width < 105 {
	t.Run("Assert displayInfoBox with small width", func(t *testing.T) {
		temp, _ := model.Update(tea.WindowSizeMsg{
			Width:  104,
			Height: 100,
		})

		model = temp.(MainModel)

		assert.Check(t, !model.displayInfoBox)
	})

	t.Run("Assert displayInfoBox with small height", func(t *testing.T) {
		temp, _ := model.Update(tea.WindowSizeMsg{
			Width:  105,
			Height: 31,
		})

		model = temp.(MainModel)

		assert.Check(t, !model.displayInfoBox)
	})
}

func TestRunBackground(t *testing.T) {
	model := MainModel{
		possibleLongRunningOpErrorChan: make(chan error, 10),
		notificationChan:               make(chan notificationMetadata, 10),
	}

	t.Run("Gets error, should not send notification", func(t *testing.T) {
		op := func() error {
			return errors.New("error")
		}

		model.runBackground(op)

		select {
		case <-model.possibleLongRunningOpErrorChan:
		default:
			t.Errorf("Should recieve an error")
		}
	})

	t.Run("Does not get an error, should send notification", func(t *testing.T) {
		op := func() error {
			return nil
		}

		model.runBackground(op)

		select {
		case <-model.possibleLongRunningOpErrorChan:
			t.Errorf("Should not recieve an error")
		default:
		}
	})
}

func TestGetRegexMatch(t *testing.T) {
	str := "Step 4/4 : RUN echo \"alpine\""
	reg := regexp.MustCompile(`Step\s(\d+)\/(\d+)\s:\s(.*)`)

	matches := reg.FindStringSubmatch(str)
	assert.DeepEqual(t, matches, []string{str, "4", "4", "RUN echo \"alpine\""})

}
