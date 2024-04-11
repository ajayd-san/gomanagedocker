package tui

import (
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
)

type dockerRes interface {
	list.Item
	getId() string
	getSize() float64
	getLabel() string
	getName() string
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

// TODO: either use this or omit this
func (i imageItem) getLabel() string {
	return "image labels here"
}

func (i imageItem) getName() string {
	return strings.Join(i.RepoTags, ", ")
}

// INFO: impl list.Item Interface
func (i imageItem) Title() string       { return i.getName() }
func (i imageItem) Description() string { return strconv.FormatFloat(i.getSize(), 'f', 2, 64) }
func (i imageItem) FilterValue() string { return i.getName() }

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

func (c containerItem) getName() string {
	return strings.Join(c.Names, ", ")
}

// INFO: impl list.Item Interface
func (i containerItem) Title() string       { return i.getName() }
func (i containerItem) Description() string { return strconv.FormatFloat(i.getSize(), 'f', 2, 64) }
func (i containerItem) FilterValue() string { return i.getLabel() }

type VolumeItem struct {
	volume.Volume
}

// FilterValue implements dockerRes.
func (v VolumeItem) FilterValue() string {
	panic("unimplemented")
}

// getId implements dockerRes.
func (v VolumeItem) getId() string {
	panic("unimplemented")
}

// getLabel implements dockerRes.
func (v VolumeItem) getLabel() string {
	panic("unimplemented")
}

// getName implements dockerRes.
func (v VolumeItem) getName() string {
	return v.Name
}

// getSize implements dockerRes.
func (v VolumeItem) getSize() float64 {
	panic("unimplemented")
}

func (i VolumeItem) Title() string       { return i.getName() }
func (i VolumeItem) Description() string { return "" }

func makeVolumeItem(dockerlist []*volume.Volume) []dockerRes {
	res := make([]dockerRes, len(dockerlist))

	for i, volume := range dockerlist {
		res[i] = VolumeItem{Volume: *volume}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].getName() < res[j].getName()
	})

	return res
}
