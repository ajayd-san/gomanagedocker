package tui

import (
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
)

type dockerRes interface {
	list.Item
	getId() string
	getSize() float64
	getLabel() string
}

type imageItem struct {
	image.Summary
}

func makeImageItems(dockerlist []image.Summary) []dockerRes {
	res := make([]dockerRes, len(dockerlist))

	for i := range dockerlist {
		res[i] = imageItem{Summary: dockerlist[i]}
	}

	return res
}

// INFO: impl dockerRes Interface
func (i imageItem) getId() string {
	return i.ID
}

func (i imageItem) getSize() float64 {
	return float64(i.Size) / float64(1e+9)
}

func (i imageItem) getLabel() string {
	return "image labels here"
}

// INFO: impl list.Item Interface
func (i imageItem) Title() string       { return i.getId() }
func (i imageItem) Description() string { return strconv.FormatFloat(i.getSize(), 'f', 2, 64) }
func (i imageItem) FilterValue() string { return i.getLabel() }

type containerItem struct {
	types.Container
}

func makeContainerItems(dockerlist []types.Container) []dockerRes {
	res := make([]dockerRes, len(dockerlist))

	for i := range dockerlist {
		res[i] = containerItem{Container: dockerlist[i]}
	}

	return res
}

// INFO: impl dockerRes Interface
func (c containerItem) getId() string {
	return c.ID
}

func (c containerItem) getSize() float64 {
	return float64(c.SizeRw) / float64(1e+9)
}

func (c containerItem) getLabel() string {
	return "labels here"
}

// INFO: impl list.Item Interface
func (i containerItem) Title() string       { return i.getId() }
func (i containerItem) Description() string { return strconv.FormatFloat(i.getSize(), 'f', 2, 64) }
func (i containerItem) FilterValue() string { return i.getLabel() }
