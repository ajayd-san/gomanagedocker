package tui

import (
	"strings"
	"testing"

	"github.com/ajayd-san/gomanagedocker/service/types"
	it "github.com/ajayd-san/gomanagedocker/service/types"
	tea "github.com/charmbracelet/bubbletea"
	"gotest.tools/v3/assert"
)

func TestUpdateExistingIds(t *testing.T) {

	containers := []types.ContainerSummary{
		{
			Names: []string{"a"},
			ID:    "1",
			Size: &it.SizeInfo{
				Rw:     1e+9,
				RootFs: 2e+9,
			},
			State: "running",
		},
		{
			Names: []string{"b"},
			ID:    "2",
			Size: &it.SizeInfo{
				Rw:     201,
				RootFs: 401,
			},
			State: "running",
		},
		{
			Names: []string{"c"},
			ID:    "3",
			Size: &it.SizeInfo{
				Rw:     202,
				RootFs: 403,
			},
			State: "running",
		},
		{

			Names: []string{"d"},
			ID:    "4",
			Size: &it.SizeInfo{
				Rw:     203,
				RootFs: 403,
			},
			State: "running",
		},
	}

	imgs := []it.ImageSummary{
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
	CONTAINERS = 0
	IMAGES = 1

	keymap := NewKeyMap(it.Docker)

	t.Run("Assert container Ids", func(t *testing.T) {
		contList := InitList(0, keymap.container, keymap.containerBulk)
		dres := makeContainerItems(containers, make(map[string]string))
		contList.updateExistigIds(&dres)
		want := map[string]struct{}{
			"1": {},
			"2": {},
			"3": {},
			"4": {},
		}

		assert.DeepEqual(t, contList.ExistingIds, want)
	})

	t.Run("Assert image Ids", func(t *testing.T) {
		imgsList := InitList(IMAGES, keymap.image, keymap.imageBulk)
		dres := makeImageItems(imgs)
		imgsList.updateExistigIds(&dres)
		want := map[string]struct{}{
			"0": {},
			"1": {},
			"2": {},
			"3": {},
		}
		assert.DeepEqual(t, imgsList.ExistingIds, want)
	})
}

func TestUpdateTab(t *testing.T) {
	IMAGES = 0
	CONTAINERS = 1

	imgs := []it.ImageSummary{
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

	keymap := NewKeyMap(it.Docker)

	list := InitList(IMAGES, keymap.image, keymap.imageBulk)
	t.Run("Assert Images subset", func(t *testing.T) {
		subset := imgs[:2]
		dres := makeImageItems(subset)
		list.updateTab(dres)

		liItems := list.list.Items()

		for i := range len(liItems) {
			got := liItems[i].(imageItem)
			want := subset[i]

			assert.DeepEqual(t, got.ImageSummary, want)
		}
	})

	t.Run("Assert Images full", func(t *testing.T) {
		dres := makeImageItems(imgs)
		list.updateTab(dres)

		liItems := list.list.Items()

		for i := range len(liItems) {
			got := liItems[i].(imageItem)
			want := imgs[i]

			assert.DeepEqual(t, got.ImageSummary, want)
		}
	})

}

func TestUpdate(t *testing.T) {
	IMAGES = 0
	CONTAINERS = 1

	imgs := []it.ImageSummary{
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

	keymap := NewKeyMap(it.Docker)
	imgList := InitList(IMAGES, keymap.image, keymap.imageBulk)

	t.Run("Update images", func(t *testing.T) {
		dres := makeImageItems(imgs)
		temp, _ := imgList.Update(dres)
		imgList = temp.(listModel)

		listItems := imgList.list.Items()

		for i := range len(listItems) {
			got := listItems[i].(imageItem)
			want := imgs[i]

			assert.DeepEqual(t, got.ImageSummary, want)
		}
	})

	t.Run("Update list size", func(t *testing.T) {
		assert.Equal(t, imgList.list.Width(), 60)
		temp, _ := imgList.Update(tea.WindowSizeMsg{Width: 210, Height: 100})
		imgList := temp.(listModel)

		assert.Equal(t, imgList.list.Width(), int(210*0.3))
	})
}

func TestEmptyList(t *testing.T) {

	IMAGES = 0
	CONTAINERS = 1

	imgs := []it.ImageSummary{
		{
			Containers: 0,
			ID:         "0as;dkfjasdfasdfasdfaasdf",
			RepoTags:   []string{"a"},
		},

		{
			Containers: 0,
			ID:         "10as;dkfjasdfasdfasdfaasdf",
			RepoTags:   []string{"b"},
		},
		{
			Containers: 3,
			ID:         "20as;dkfjasdfasdfasdfaasdf",
			RepoTags:   []string{"c"},
		},
		{
			Containers: 0,
			ID:         "30as;dkfjasdfasdfasdfaasdf",
			RepoTags:   []string{"d"},
		},
	}

	keymap := NewKeyMap(it.Docker)
	imgList := InitList(IMAGES, keymap.image, keymap.imageBulk)

	t.Run("List with items", func(t *testing.T) {
		dres := makeImageItems(imgs)
		temp, _ := imgList.Update(dres)
		imgList = temp.(listModel)
		got := imgList.View()

		assert.Assert(t, !strings.Contains(got, "No items"))
	})

	t.Run("Empty list", func(t *testing.T) {
		temp, _ := imgList.Update([]dockerRes{})
		imgList = temp.(listModel)
		got := imgList.View()

		assert.Assert(t, strings.Contains(got, "No items"))

	})

}
