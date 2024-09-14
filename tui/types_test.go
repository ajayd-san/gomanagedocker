package tui

import (
	"strings"
	"testing"

	"github.com/ajayd-san/gomanagedocker/service/types"
	podmanTypes "github.com/containers/podman/v5/pkg/domain/entities/types"
	"gotest.tools/v3/assert"
)

func TestMakeDescriptionString(t *testing.T) {
	str1, str2 := "mad", "scientist"

	listContainerWidth := listContainer.GetWidth()
	got := makeDescriptionString(str1, str2, len(str1))
	want := str1 + strings.Repeat(" ", listContainerWidth-len(str1)-len(str2)-3) + str2

	assert.Equal(t, got, want)
}

func TestMakeImageItems(t *testing.T) {
	dockerList := []types.ImageSummary{
		{
			ID:       "#1",
			RepoTags: []string{"latest", "tag1", "tag2"},
		}, {
			ID:       "#2",
			RepoTags: []string{},
		},
	}

	t.Run("Should only return non dangling items", func(t *testing.T) {
		// for dangling items the repo tags would be an empty slice
		got := makeImageItems(dockerList)
		assert.Equal(t, len(got), 1)

		for _, item := range got {
			imgItem, ok := item.(imageItem)
			assert.Equal(t, ok, true)
			assert.Assert(t, len(imgItem.RepoTags) != 0)
		}
	})

}

func TestTransformListNames(t *testing.T) {

	names := []string{"name1", "name2", "name3", "name4"}

	t.Run("With listContainer.Width = 13", func(t *testing.T) {
		listContainer = listContainer.Width(13)
		got := transformListNames(names)
		want := "name1"
		assert.Equal(t, got, want)
	})

	t.Run("With listContainer.Width = 12", func(t *testing.T) {
		listContainer = listContainer.Width(12)
		got := transformListNames(names)
		want := "name1"
		assert.Equal(t, got, want)
	})

	t.Run("With listContainer.Width = 20", func(t *testing.T) {
		listContainer = listContainer.Width(20)
		got := transformListNames(names)
		want := "name1, name2"
		assert.Equal(t, got, want)
	})

	t.Run("With listContainer.Width=56", func(t *testing.T) {
		listContainer = listContainer.Width(56)
		names := []string{"a.smol.list:latest", "alpine:latest", "b.star.man:latest"}
		got := transformListNames(names)
		want := "a.smol.list:latest, alpine:latest"
		assert.Equal(t, got, want)
	})

	t.Run("with listContainer.Width=20, Edge case", func(t *testing.T) {
		listContainer = listContainer.Width(20)
		names := []string{"Zenitsu", "best"}
		got := transformListNames(names)
		want := "Zenitsu, best"
		assert.Equal(t, got, want)
	})

	t.Run("With empty list", func(t *testing.T) {
		defer func() {
			if recover() != nil {
				t.Error("This function should not panic")
			}
		}()
		listContainer = listContainer.Width(20)
		names := make([]string, 0)
		// should not panic
		transformListNames(names)
	})
}

func TestGetRunningContainers(t *testing.T) {
	item := PodItem{
		ListPodsReport: podmanTypes.ListPodsReport{
			Containers: []*podmanTypes.ListPodContainer{
				{
					Id:     "a",
					Status: "running",
				},
				{
					Id:     "b",
					Status: "running",
				},
				{
					Id:     "c",
					Status: "running",
				},
				{
					Id:     "d",
					Status: "exited",
				},
			},
			Id:     "1234",
			Name:   "mario",
			Status: "running",
		},
	}

	got := item.getRunningContainers()
	want := 3

	assert.Equal(t, got, want)
}
