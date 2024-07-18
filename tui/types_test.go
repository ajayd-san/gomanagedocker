package tui

import (
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestMakeDescriptionString(t *testing.T) {
	str1, str2 := "mad", "scientist"

	listContainerWidth := listContainer.GetWidth()
	got := makeDescriptionString(str1, str2, len(str1))
	want := str1 + strings.Repeat(" ", listContainerWidth-len(str1)-len(str2)-3) + str2

	assert.Equal(t, got, want)
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
		listContainer = listContainer.Width(20)
		names := make([]string, 0)
		// will panic
		transformListNames(names)
	})
}
