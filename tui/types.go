package tui

import (
	"cmp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
)

type status int

const (
	containerStateRunning status = iota
	containerStatePaused
	containerStateRestarting
	containerStateExited
	containerStateCreated
	containerStateRemoving
	containerStateDead
)

var statusMap = map[string]status{
	"running":    containerStateRunning,
	"paused":     containerStatePaused,
	"restarting": containerStateRestarting,
	"exited":     containerStateExited,
	"created":    containerStateCreated,
	"removing":   containerStateRemoving,
	"dead":       containerStateDead,
}

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
	res := make([]dockerRes, 0)

	for i := range dockerlist {
		if len(dockerlist[i].RepoTags) == 0 {
			continue
		}

		res = append(res, imageItem{Summary: dockerlist[i]})
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
	return transformListNames(i.RepoTags)
}

// INFO: impl list.Item Interface
func (i imageItem) Title() string { return i.getName() }

func (i imageItem) Description() string {
	id := i.getId()
	id = strings.TrimPrefix(id, "sha256:")
	shortId := id[:15]

	sizeStr := strconv.FormatFloat(i.getSize(), 'f', 2, 64) + "GB"

	return makeDescriptionString(shortId, sizeStr, len(shortId))
}

func (i imageItem) FilterValue() string { return i.getName() }

type containerItem struct {
	types.Container
}

func makeContainerItems(dockerlist []types.Container) []dockerRes {
	res := make([]dockerRes, len(dockerlist))

	slices.SortFunc(dockerlist, func(a types.Container, b types.Container) int {

		if statusMap[a.State] < statusMap[b.State] {
			return -1
		} else if statusMap[a.State] > statusMap[b.State] {
			return 1
		}

		// we can compare by only first name, since names cannot be equal
		return cmp.Compare(a.Names[0], b.Names[0])
	})

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
	return c.getName()
}

func (c containerItem) getName() string {
	return transformListNames(c.Names)
}

func (c containerItem) getState() string {
	return c.State
}

// INFO: impl list.Item Interface
func (i containerItem) Title() string { return i.getName() }

func (i containerItem) Description() string {

	id := i.getId()
	id = strings.TrimPrefix(id, "sha256:")
	shortId := id[:15]

	state := i.State
	switch i.State {
	case "running":
		state = containerRunningStyle.Render(state)
	case "exited":
		state = containerExitedStyle.Render(state)
	case "created":
		state = containerCreatedStyle.Render(state)
	case "restarting":
		state = containerRestartingStyle.Render(state)
	case "dead":
		state = containerDeadStyle.Render(state)
	}

	return makeDescriptionString(shortId, state, len(shortId))
}

func (i containerItem) FilterValue() string { return i.getLabel() }

type VolumeItem struct {
	volume.Volume
}

func (v VolumeItem) FilterValue() string {
	return v.getId()
}

func (v VolumeItem) getId() string {
	return v.Name
}

func (v VolumeItem) getLabel() string {
	panic("unimplemented")
}

func (v VolumeItem) getName() string {

	return v.Name[:min(30, len(v.Name))]
}

func (v VolumeItem) getSize() float64 {
	if v.UsageData == nil {
		return -1
	}
	return float64(v.UsageData.Size)
}

func (i VolumeItem) Title() string { return i.getName() }

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

// util

/*
This function makes the final description string with white space between the two strings
using string manipulation, offset is typically the length of the first string.
The final length of the returned string would be listContainer.Width - offset - 3
*/
func makeDescriptionString(str1, str2 string, offset int) string {
	str2 = lipgloss.PlaceHorizontal(listContainer.GetWidth()-offset-3, lipgloss.Right, str2)
	return lipgloss.JoinHorizontal(lipgloss.Left, str1, str2)
}

// This function takes in names associated with objects (e.g: RepoTags in case of Image)
// and concatenates into a string depending on the width of the list
func transformListNames(names []string) string {
	if len(names) == 0 {
		return ""
	}

	runningLength := 0
	var maxindex int
	for index, name := range names {
		runningLength += len(name)
		if runningLength > listContainer.GetWidth()-7 {
			break
		}
		if index != len(names)-1 {
			runningLength += 2 // +2 cuz we also append ", " after each element
		}
		maxindex = index
	}

	res := strings.Join(names[:maxindex+1], ", ")

	if len(res) > listContainer.GetWidth()-7 {
		return res[:listContainer.GetWidth()-7] + "..."
	}

	return res
}
