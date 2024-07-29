package tui

import (
	"errors"
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

	wg := sync.WaitGroup{}
	newlist := model.fetchNewData(0, &wg)
	wg.Wait()

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
			InitList(0),
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
			InitList(0),
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

		notif := NewNotification(1, "Test notification")
		model.runBackground(op, &notif)

		select {
		case <-model.possibleLongRunningOpErrorChan:
		default:
			t.Errorf("Should recieve an error")
		}

		select {
		case <-model.notificationChan:
			t.Errorf("Should not recieve a notification")
		default:
		}
	})

	t.Run("Does not get an error, should send notification", func(t *testing.T) {
		op := func() error {
			return nil
		}

		notif := NewNotification(1, "Test notification")
		model.runBackground(op, &notif)

		select {
		case <-model.possibleLongRunningOpErrorChan:
			t.Errorf("Should not recieve an error")
		default:
		}

		select {
		case <-model.notificationChan:
		default:
			t.Errorf("Should recieve a notification")
		}
	})
}
