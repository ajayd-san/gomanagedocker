package tui

import (
	"log"
	"slices"
	"testing"

	"github.com/ajayd-san/gomanagedocker/dockercmd"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"golang.design/x/clipboard"
	"gotest.tools/v3/assert"
)

func setupTest(t *testing.T) dockercmd.DockerClient {
	api := dockercmd.MockApi{}

	containers := []types.Container{
		{
			Names:      []string{"a"},
			ID:         "1aaaaaaaa",
			SizeRw:     1e+9,
			SizeRootFs: 2e+9,
			State:      "running",
			Status:     "",
		},
		{
			Names:      []string{"b"},
			ID:         "2aaaaaaaa",
			SizeRw:     201,
			SizeRootFs: 401,
			State:      "running",
		},
		{
			Names:      []string{"c"},
			ID:         "3aaaaaaaa",
			SizeRw:     202,
			SizeRootFs: 402,
			State:      "running",
		},
		{
			Names:      []string{"d"},
			ID:         "4aaaaaaaa",
			SizeRw:     203,
			SizeRootFs: 403,
			State:      "running",
		},
	}

	imgs := []image.Summary{
		{
			Containers: 0,
			ID:         "0bbbbbbbb",
			RepoTags:   []string{"a"},
		},

		{
			Containers: 0,
			ID:         "1bbbbbbbb",
			RepoTags:   []string{"b"},
		},
		{
			Containers: 3,
			ID:         "2bbbbbbbb",
			RepoTags:   []string{"c"},
		},
		{
			Containers: 0,
			ID:         "3bbbbbbbb",
			RepoTags:   []string{"d"},
		},
	}

	api.SetMockContainers(containers)
	api.SetMockImages(imgs)

	mock := dockercmd.NewMockCli(&api)
	return mock
}

func TestToggleStartStopContainer(t *testing.T) {

	tests := []struct {
		target    types.Container
		want      string
		notifWant string
	}{
		{
			target: types.Container{
				Names:      []string{"b"},
				ID:         "2aaaaaaaa",
				SizeRw:     201,
				SizeRootFs: 401,
				State:      "running",
			},
			want:      "stopped",
			notifWant: listStatusMessageStyle.Render("Stopped 2aaaaaaa"),
		},
		{
			target: types.Container{
				Names:      []string{"b"},
				ID:         "2aaaaaaaa",
				SizeRw:     201,
				SizeRootFs: 401,
				State:      "stopped",
			},
			want:      "running",
			notifWant: listStatusMessageStyle.Render("Started 2aaaaaaa"),
		},
	}

	mock := setupTest(t)
	mock.ToggleContainerListAll()

	for _, testCase := range tests {
		t.Run("Test for existing container", func(t *testing.T) {
			target := containerItem{
				testCase.target,
			}

			notifChan := make(chan notificationMetadata, 10)
			op := toggleStartStopContainer(mock, target, 1, notifChan)

			op()

			t.Run("Test Stopping", func(t *testing.T) {
				containers := mock.ListContainers(false)
				index := slices.IndexFunc(containers, func(elem types.Container) bool {
					if elem.ID == target.ID {
						return true
					}

					return false
				})

				got := containers[index]

				assert.Equal(t, got.State, testCase.want)
			})

			t.Run("Assert Notification", func(t *testing.T) {
				select {
				case notif := <-notifChan:
					assert.Equal(t, notif, notificationMetadata{
						listId: 1,
						msg:    testCase.notifWant,
					})
				default:
					t.Errorf("No notification received")
				}
			})
		})
	}
}

func TestTogglePauseResumeContainer(t *testing.T) {
	mock := setupTest(t)

	tests := []struct {
		target    types.Container
		want      string
		notifWant string
	}{
		{
			target: types.Container{
				Names:      []string{"b"},
				ID:         "2aaaaaaaa",
				SizeRw:     201,
				SizeRootFs: 401,
				State:      "running",
			},
			want:      "paused",
			notifWant: listStatusMessageStyle.Render("Paused 2aaaaaaa"),
		},
		{
			target: types.Container{
				Names:      []string{"b"},
				ID:         "2aaaaaaaa",
				SizeRw:     201,
				SizeRootFs: 401,
				State:      "paused",
			},
			want:      "running",
			notifWant: listStatusMessageStyle.Render("Resumed 2aaaaaaa"),
		},
	}

	for _, testCase := range tests {

		t.Run("Test for Existing Container", func(t *testing.T) {

			target := containerItem{
				testCase.target,
			}

			notifChan := make(chan notificationMetadata, 10)
			op := togglePauseResumeContainer(mock, target, 2, notifChan)

			op()

			t.Run("Assert Paused State", func(t *testing.T) {
				containers := mock.ListContainers(false)
				log.Println(containers)

				index := slices.IndexFunc(containers, func(elem types.Container) bool {
					if elem.ID == target.ID {
						return true
					}

					return false
				})

				got := containers[index]

				assert.Equal(t, got.State, testCase.want)
			})

			t.Run("Assert Notification", func(t *testing.T) {
				select {
				case notif := <-notifChan:
					assert.Equal(t, notif, notificationMetadata{
						listId: 2,
						msg:    testCase.notifWant,
					})
				default:
					t.Errorf("No notification received")
				}
			})

		})

	}
}

func TestContainerDeleteForce(t *testing.T) {
	tests := []struct {
		target    types.Container
		notifWant string
		errorStr  string
	}{
		{
			target: types.Container{
				Names:      []string{"b"},
				ID:         "2aaaaaaaa",
				SizeRw:     201,
				SizeRootFs: 401,
				State:      "running",
			},
			notifWant: listStatusMessageStyle.Render("Deleted 2aaaaaaa"),
		},
		{
			target: types.Container{
				Names:      []string{"xyz"},
				ID:         "this container does not exist",
				SizeRw:     201,
				SizeRootFs: 401,
				State:      "running",
			},
			notifWant: "",
			errorStr:  "No such container:",
		},
	}

	mock := setupTest(t)
	mock.ToggleContainerListAll()

	for _, testCase := range tests {
		t.Run("Force Delete Exising Container", func(t *testing.T) {
			target := containerItem{testCase.target}

			notifChan := make(chan notificationMetadata, 10)
			op := containerDeleteForce(mock, target, 2, notifChan)

			err := op()

			// test for error
			if testCase.errorStr != "" {
				assert.ErrorContains(t, err, testCase.errorStr)
				// if there is an error, return early so that we do not perform other subtests
				return
			}

			t.Run("Confirm container deleted", func(t *testing.T) {
				containers := mock.ListContainers(false)

				exists := slices.ContainsFunc(containers, func(elem types.Container) bool {
					if elem.ID == target.ID {
						return true
					}
					return false
				})

				assert.Assert(t, !exists)
			})

			t.Run("Assert Notification", func(t *testing.T) {
				select {
				case notif := <-notifChan:
					assert.Equal(t, notif, notificationMetadata{
						listId: 2,
						msg:    testCase.notifWant,
					})
				default:
					t.Errorf("No notification received")
				}
			})
		})
	}
}

func TestCopyIdToClipboard(t *testing.T) {
	clipboard.Init()
	target := containerItem{
		types.Container{
			Names:      []string{"b"},
			ID:         "TuTuRuu!",
			SizeRw:     201,
			SizeRootFs: 401,
			State:      "running",
		},
	}

	notifChan := make(chan notificationMetadata, 10)
	op := copyIdToClipboard(target, 1, notifChan)
	op()

	got := clipboard.Read(clipboard.FmtText)
	assert.Equal(t, string(got), target.ID)
}

func TestImageDeleteForce(t *testing.T) {
	tests := []struct {
		target    image.Summary
		notifWant string
		errorStr  string
	}{
		{
			target: image.Summary{
				Containers: 0,
				ID:         "0bbbbbbbb",
				RepoTags:   []string{"a"},
			},

			notifWant: listStatusMessageStyle.Render("Deleted 0bbbbbbb"),
		},
		{
			target: image.Summary{
				Containers: 0,
				ID:         "0bbbbbbbb",
				RepoTags:   []string{"a"},
			},
			notifWant: "",
			errorStr:  "No such image:",
		},
	}

	mock := setupTest(t)
	mock.ToggleContainerListAll()

	for _, testCase := range tests {
		t.Run("Force Delete Exising image", func(t *testing.T) {
			target := imageItem{testCase.target}

			notifChan := make(chan notificationMetadata, 10)
			op := imageDeleteForce(mock, target, 2, notifChan)

			err := op()

			// test for error
			if testCase.errorStr != "" {
				assert.ErrorContains(t, err, testCase.errorStr)
				// if there is an error, return early so that we do not perform other subtests
				return
			}

			t.Run("Confirm image deleted", func(t *testing.T) {
				images := mock.ListImages()

				exists := slices.ContainsFunc(images, func(elem image.Summary) bool {
					if elem.ID == target.ID {
						return true
					}
					return false
				})

				assert.Assert(t, !exists)
			})

			t.Run("Assert Notification", func(t *testing.T) {
				select {
				case notif := <-notifChan:
					assert.Equal(t, notif, notificationMetadata{
						listId: 2,
						msg:    testCase.notifWant,
					})
				default:
					t.Errorf("No notification received")
				}
			})
		})
	}
}
