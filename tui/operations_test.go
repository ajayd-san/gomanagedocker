package tui

import (
	"slices"
	"testing"

	"github.com/ajayd-san/gomanagedocker/service"
	"github.com/ajayd-san/gomanagedocker/service/dockercmd"
	"github.com/ajayd-san/gomanagedocker/service/podmancmd"
	it "github.com/ajayd-san/gomanagedocker/service/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
	"gotest.tools/v3/assert"

	// podman types

	pt "github.com/containers/podman/v5/pkg/domain/entities/types"
	defineSize "github.com/containers/podman/v5/pkg/ps/define"

	"github.com/containers/podman/v5/libpod/define"
)

func setupMockDockerClient(t *testing.T) service.Service {
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
			State:      "stopped",
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

	vols := []*volume.Volume{
		{
			Name: "1",
		},
		{
			Name: "2",
		},
		{
			Name: "3",
		},
	}

	api.SetMockContainers(containers)
	api.SetMockImages(imgs)
	api.SetMockVolumes(vols)

	mock := dockercmd.NewMockCli(&api)
	return mock
}
func setupMockPodmanClient(t *testing.T) service.Service {
	api := podmancmd.PodmanMockApi{}

	containers := []pt.ListContainer{
		{
			Names: []string{"a"},
			ID:    "1aaaaaaaa",
			Size: &defineSize.ContainerSize{
				RootFsSize: 2e+9,
				RwSize:     1e+9,
			},
			State:  "running",
			Status: "",
		},
		{
			Names: []string{"b"},
			ID:    "2aaaaaaaa",
			Size: &defineSize.ContainerSize{
				RootFsSize: 201,
				RwSize:     401,
			},
			State: "running",
		},
		{
			Names: []string{"c"},
			ID:    "3aaaaaaaa",
			Size: &defineSize.ContainerSize{
				RootFsSize: 202,
				RwSize:     402,
			},
			State: "running",
		},
		{
			Names: []string{"d"},
			ID:    "4aaaaaaaa",
			Size: &defineSize.ContainerSize{
				RootFsSize: 203,
				RwSize:     403,
			},
			State: "stopped",
		},
	}

	imgs := []*pt.ImageSummary{
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

	vols := []*pt.VolumeListReport{
		{
			VolumeConfigResponse: pt.VolumeConfigResponse{
				InspectVolumeData: define.InspectVolumeData{
					Name: "1",
				},
			},
		},
		{
			VolumeConfigResponse: pt.VolumeConfigResponse{
				InspectVolumeData: define.InspectVolumeData{
					Name: "2",
				},
			},
		},
		{
			VolumeConfigResponse: pt.VolumeConfigResponse{
				InspectVolumeData: define.InspectVolumeData{
					Name: "3",
				},
			},
		},
	}

	pods := []*pt.ListPodsReport{
		{
			Id:     "1",
			Name:   "ippo",
			Status: "running",
		},
		{
			Id:     "2",
			Name:   "zenitsu",
			Status: "running",
		},
		{
			Id:     "3",
			Name:   "gojo",
			Status: "stopped",
		},
	}

	api.SetMockContainers(containers)
	api.SetMockImages(imgs)
	api.SetMockVolumes(vols)
	api.SetMockPods(pods)

	mock := podmancmd.NewMockCli(&api)
	return mock
}

func TestToggleStartStopContainer(t *testing.T) {

	tests := []struct {
		containers []dockerRes
		wantState  []string
		notifWant  []string
	}{
		{
			containers: []dockerRes{
				containerItem{
					it.ContainerSummary{
						Names: []string{"a"},
						ID:    "1aaaaaaaa",
						Size: &it.SizeInfo{
							Rw:     1e+9,
							RootFs: 2e+9,
						},
						State: "running",
					},
					"abc",
				},
				containerItem{
					it.ContainerSummary{
						Names: []string{"b"},
						ID:    "2aaaaaaaa",
						Size: &it.SizeInfo{
							Rw:     201,
							RootFs: 401,
						},
						State: "running",
					},
					"idk",
				},
			},
			wantState: []string{"stopped", "stopped"},
			notifWant: []string{
				listStatusMessageStyle.Render("Stopped 1aaaaaaa"),
				listStatusMessageStyle.Render("Stopped 2aaaaaaa"),
				listStatusMessageStyle.Render("Toggled 2 containers"),
			},
		},
		{
			containers: []dockerRes{
				containerItem{
					it.ContainerSummary{
						Names: []string{"b"},
						ID:    "2aaaaaaaa",
						Size: &it.SizeInfo{
							Rw:     201,
							RootFs: 401,
						},
						State: "stopped",
					},

					"idk",
				},
				containerItem{
					it.ContainerSummary{
						Names: []string{"b"},
						ID:    "3aaaaaaaa",
						Size: &it.SizeInfo{
							Rw:     201,
							RootFs: 401,
						},
						State: "running",
					},

					"idk",
				},
			},
			wantState: []string{"running", "stopped"},
			notifWant: []string{
				listStatusMessageStyle.Render("Started 2aaaaaaa"),
				listStatusMessageStyle.Render("Started 3aaaaaaa"),
				listStatusMessageStyle.Render("Toggled 2 containers"),
			},
		},
	}

	test := func(t *testing.T, mock service.Service) {
		for _, testCase := range tests {
			t.Run("Test for existing container", func(t *testing.T) {

				notifChan := make(chan notificationMetadata, 10)
				errChan := make(chan error, 10)
				op := toggleStartStopContainer(mock, testCase.containers, 1, notifChan, errChan)

				op()

				t.Run("Test Stopping", func(t *testing.T) {
					updatedContainers := mock.ListContainers(false)
					for i, container := range testCase.containers {
						id := container.GetId()

						index := slices.IndexFunc(updatedContainers, func(elem it.ContainerSummary) bool {
							return elem.ID == id
						})

						assert.Equal(t, updatedContainers[index].State, testCase.wantState[i])
					}
				})

				t.Run("Assert Notification", func(t *testing.T) {
					// ik this is not a complete test but it's just easier.
					// TODO: assert each notification
					assert.Equal(t, len(testCase.notifWant), len(notifChan))
				})
			})
		}
	}

	t.Run("docker", func(t *testing.T) {
		mock := setupMockDockerClient(t)
		mock.ToggleContainerListAll()

		test(t, mock)
	})

	t.Run("podman", func(t *testing.T) {
		mock := setupMockPodmanClient(t)
		mock.ToggleContainerListAll()

		test(t, mock)
	})
}

func TestTogglePauseResumeContainer(t *testing.T) {
	tests := []struct {
		containers []dockerRes
		wantState  []string
		notifWant  []string
	}{
		{
			containers: []dockerRes{
				containerItem{
					it.ContainerSummary{
						Names: []string{"a"},
						ID:    "1aaaaaaaa",
						Size: &it.SizeInfo{
							Rw:     1e+9,
							RootFs: 2e+9,
						},
						State: "running",
					},
					"abc",
				},
				containerItem{
					it.ContainerSummary{
						Names: []string{"b"},
						ID:    "2aaaaaaaa",
						Size: &it.SizeInfo{
							Rw:     201,
							RootFs: 401,
						},
						State: "running",
					},
					"idk",
				},
			},

			wantState: []string{"paused", "paused"},
			notifWant: []string{
				listStatusMessageStyle.Render("Paused 1aaaaaaa"),
				listStatusMessageStyle.Render("Paused 2aaaaaaa"),
				listStatusMessageStyle.Render("Paused: 2 containers"),
			},
		},
		{
			containers: []dockerRes{
				containerItem{
					it.ContainerSummary{
						Names: []string{"b"},
						ID:    "2aaaaaaaa",
						Size: &it.SizeInfo{
							Rw:     201,
							RootFs: 401,
						},
						State: "paused",
					},
					"idk",
				},
				containerItem{
					it.ContainerSummary{
						Names: []string{"b"},
						ID:    "3aaaaaaaa",
						Size: &it.SizeInfo{
							Rw:     201,
							RootFs: 401,
						},
						State: "running",
					},
					"idk",
				},
			},
			wantState: []string{"running", "paused"},
			notifWant: []string{
				listStatusMessageStyle.Render("Resumed 2aaaaaaa"),
				listStatusMessageStyle.Render("Paused 3aaaaaaa"),
				listStatusMessageStyle.Render("Paused: 1, Resumed: 1 containers"),
			},
		},
	}

	test := func(t *testing.T, mock service.Service) {
		for _, testCase := range tests {
			t.Run("Test for Existing Container", func(t *testing.T) {
				notifChan := make(chan notificationMetadata, 10)
				errChan := make(chan error, 10)
				op := togglePauseResumeContainer(mock, testCase.containers, 2, notifChan, errChan)

				op()

				t.Run("Assert Paused State", func(t *testing.T) {
					updatedContainers := mock.ListContainers(false)
					for i, container := range testCase.containers {
						id := container.GetId()

						index := slices.IndexFunc(updatedContainers, func(elem it.ContainerSummary) bool {
							return elem.ID == id
						})

						assert.Equal(t, updatedContainers[index].State, testCase.wantState[i])
					}
				})

				t.Run("Assert Notification", func(t *testing.T) {
					assert.Equal(t, len(testCase.notifWant), len(notifChan))
					// select {
					// case notif := <-notifChan:
					// 	assert.Equal(t, notif, notificationMetadata{
					// 		listId: 2,
					// 		msg:    testCase.notifWant,
					// 	})
					// default:
					// 	t.Errorf("No notification received")
					// }
				})

			})

		}
	}

	t.Run("docker", func(t *testing.T) {
		mock := setupMockDockerClient(t)
		test(t, mock)

	})
	t.Run("podman", func(t *testing.T) {
		mock := setupMockPodmanClient(t)
		test(t, mock)
	})

}

func TestContainerDeleteBulk(t *testing.T) {
	tests := []struct {
		containers []dockerRes
		notifWant  []string
	}{
		{
			containers: []dockerRes{
				containerItem{
					it.ContainerSummary{
						Names: []string{"a"},
						ID:    "1aaaaaaaa",
						Size: &it.SizeInfo{
							Rw:     1e+9,
							RootFs: 2e+9,
						},
						State: "running",
					},
					"abc",
				},
				containerItem{
					it.ContainerSummary{
						Names: []string{"b"},
						ID:    "2aaaaaaaa",
						Size: &it.SizeInfo{
							Rw:     201,
							RootFs: 401,
						},
						State: "running",
					},
					"idk",
				},
			},
			notifWant: []string{
				listStatusMessageStyle.Render("Deleted 1aaaaaaa"),
				listStatusMessageStyle.Render("Deleted 2aaaaaaa"),
				listStatusMessageStyle.Render("Deleted 2 containers"),
			},
		},
		{
			containers: []dockerRes{
				containerItem{
					it.ContainerSummary{
						Names: []string{"b"},
						ID:    "2aaaaaaaa",
						Size: &it.SizeInfo{
							Rw:     201,
							RootFs: 401,
						},
						State: "paused",
					},
					"idk",
				},
			},
			notifWant: []string{
				listStatusMessageStyle.Render("Deleted 2aaaaaaa"),
			},
		},
	}

	opts := it.ContainerRemoveOpts{
		RemoveVolumes: false,
		RemoveLinks:   false,
		Force:         true,
	}

	test := func(t *testing.T, mock service.Service, testCase struct {
		containers []dockerRes
		notifWant  []string
	}) {
		t.Run("Force Delete Exising Container", func(t *testing.T) {

			notifChan := make(chan notificationMetadata, 10)
			errChan := make(chan error, 10)
			op := containerDeleteBulk(mock, testCase.containers, opts, 2, notifChan, errChan)

			op()

			t.Run("Confirm container deleted", func(t *testing.T) {
				containers := mock.ListContainers(true)

				exists := slices.ContainsFunc(containers, func(elem it.ContainerSummary) bool {
					for _, c := range testCase.containers {
						if elem.ID == c.GetId() {
							return true
						}
					}
					return false
				})

				assert.Assert(t, !exists)
			})

			t.Run("Assert Notification", func(t *testing.T) {
				assert.Equal(t, len(testCase.notifWant), len(notifChan))

				for range testCase.notifWant {
					select {
					case notif := <-notifChan:
						found := slices.Contains(testCase.notifWant, notif.msg)
						assert.Assert(t, found)

					default:
						t.Errorf("No notification received")
					}
				}
			})
		})
	}

	t.Run("docker", func(t *testing.T) {
		for _, testCase := range tests {
			mock := setupMockDockerClient(t)
			mock.ToggleContainerListAll()

			test(t, mock, testCase)
		}

	})

	t.Run("podman", func(t *testing.T) {
		for _, testCase := range tests {
			mock := setupMockPodmanClient(t)
			mock.ToggleContainerListAll()
			test(t, mock, testCase)
		}
	})

}

func TestContainerDelete(t *testing.T) {
	tests := []struct {
		ID        string
		notifWant string
		errorStr  string
		opts      it.ContainerRemoveOpts
	}{
		{
			ID:        "2aaaaaaaa",
			notifWant: listStatusMessageStyle.Render("Deleted 2aaaaaaa"),
			opts: it.ContainerRemoveOpts{
				RemoveVolumes: false,
				RemoveLinks:   false,
				Force:         true,
			},
		},
		{
			ID:        "4aaaaaaaa",
			notifWant: listStatusMessageStyle.Render("Deleted 4aaaaaaa"),
		},
		{
			ID:        "3aaaaaaaa",
			notifWant: listStatusMessageStyle.Render("Deleted 3aaaaaaa"),
			errorStr:  "container is running",
		},
		{
			ID:        "this container does not exist",
			notifWant: "",
			errorStr:  "No such container:",
		},
	}

	test := func(t *testing.T, mock service.Service) {
		for _, testCase := range tests {
			t.Run("Force Delete Exising Container", func(t *testing.T) {

				notifChan := make(chan notificationMetadata, 10)
				t.Log(testCase.ID)
				op := containerDelete(mock, testCase.ID, testCase.opts, 2, notifChan)

				err := op()

				// test for error
				if testCase.errorStr != "" {
					assert.ErrorContains(t, err, testCase.errorStr)
					// if there is an error, return early so that we do not perform other subtests
					return
				}

				t.Run("Confirm container deleted", func(t *testing.T) {
					containers := mock.ListContainers(false)

					exists := slices.ContainsFunc(containers, func(elem it.ContainerSummary) bool {
						if elem.ID == testCase.ID {
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

	t.Run("docker", func(t *testing.T) {
		mock := setupMockDockerClient(t)
		mock.ToggleContainerListAll()
		test(t, mock)
	})

	t.Run("podman", func(t *testing.T) {
		mock := setupMockPodmanClient(t)
		mock.ToggleContainerListAll()
		test(t, mock)
	})

}

// this works but doesn't work on CI even with libx11-dev
// func TestCopyIdToClipboard(t *testing.T) {
// 	clipboard.Init()
// 	target := containerItem{
// 		types.Container{
// 			Names:      []string{"b"},
// 			ID:         "TuTuRuu!",
// 			SizeRw:     201,
// 			SizeRootFs: 401,
// 			State:      "running",
// 		},
// 	}

// 	notifChan := make(chan notificationMetadata, 10)
// 	op := copyIdToClipboard(target, 1, notifChan)
// 	op()

// 	got := clipboard.Read(clipboard.FmtText)
// 	assert.Equal(t, string(got), target.ID)
// }

func TestImageDelete(t *testing.T) {
	tests := []struct {
		ID        string
		notifWant string
		errorStr  string
		opts      it.RemoveImageOptions
	}{
		{
			ID:        "0bbbbbbbb",
			notifWant: listStatusMessageStyle.Render("Deleted 0bbbbbbb"),
			opts: it.RemoveImageOptions{
				Force:   false,
				NoPrune: true,
			},
		},
		{
			ID:        "0bbbbbbbb",
			notifWant: "",
			errorStr:  "No such image:",
		},
		// Should fail, since the image running with this ID has active containers assosicated
		{
			ID:       "2bbbbbbbb",
			errorStr: "unable to delete",
		},
	}

	test := func(t *testing.T, mock service.Service) {
		for _, testCase := range tests {
			t.Run("Force Delete Exising image", func(t *testing.T) {

				notifChan := make(chan notificationMetadata, 10)
				op := imageDelete(mock, testCase.ID, testCase.opts, 2, notifChan)

				err := op()

				// test for error
				if testCase.errorStr != "" {
					assert.ErrorContains(t, err, testCase.errorStr)
					// if there is an error, return early so that we do not perform other subtests
					return
				}

				t.Run("Confirm image deleted", func(t *testing.T) {
					images := mock.ListImages()

					exists := slices.ContainsFunc(images, func(elem it.ImageSummary) bool {
						if elem.ID == testCase.ID {
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

	t.Run("docker", func(t *testing.T) {
		mock := setupMockDockerClient(t)
		mock.ToggleContainerListAll()
		test(t, mock)
	})

	t.Run("podman", func(t *testing.T) {
		mock := setupMockPodmanClient(t)
		mock.ToggleContainerListAll()
		test(t, mock)
	})

}
func TestImageDeleteBulk(t *testing.T) {
	tests := []struct {
		imgs   []dockerRes
		notifs []string
		errors []string
		opts   it.RemoveImageOptions
	}{
		{
			imgs: []dockerRes{
				imageItem{
					it.ImageSummary{
						Containers: 0,
						ID:         "0bbbbbbbb",
						RepoTags:   []string{"a"},
					},
				},

				imageItem{
					it.ImageSummary{
						Containers: 0,
						ID:         "1bbbbbbbb",
						RepoTags:   []string{"b"},
					},
				},
			},
			notifs: []string{
				listStatusMessageStyle.Render("Deleted 0bbbbbbb"),
				listStatusMessageStyle.Render("Deleted 1bbbbbbb"),
				listStatusMessageStyle.Render("Deleted 2 images"),
			},
			opts: it.RemoveImageOptions{
				Force:   true,
				NoPrune: true,
			},
		},
		{
			imgs: []dockerRes{
				imageItem{
					it.ImageSummary{
						Containers: 0,
						ID:         "this does not exist",
						RepoTags:   []string{"a"},
					},
				},

				imageItem{
					it.ImageSummary{
						Containers: 0,
						ID:         "1bbbbbbbb",
						RepoTags:   []string{"b"},
					},
				},
			},
			notifs: []string{
				listStatusMessageStyle.Render("Deleted 1bbbbbbb"),
			},
			opts: it.RemoveImageOptions{
				Force:   true,
				NoPrune: true,
			},
		},
	}

	for _, testCase := range tests {
		mock := setupMockDockerClient(t)
		mock.ToggleContainerListAll()
		t.Run("Force Delete Exising image", func(t *testing.T) {

			notifChan := make(chan notificationMetadata, 10)
			erroChan := make(chan error, 10)
			op := imageDeleteBulk(mock, testCase.imgs, testCase.opts, 2, notifChan, erroChan)

			_ = op()

			// test for error
			if testCase.errors != nil {
				for range testCase.errors {
					select {
					case err := <-erroChan:
						slices.ContainsFunc(testCase.errors, func(elem string) bool {
							return err.Error() == elem
						})

					default:
					}
				}
			}

			t.Run("Confirm image deleted", func(t *testing.T) {
				images := mock.ListImages()

				exists := slices.ContainsFunc(images, func(elem it.ImageSummary) bool {
					for _, dres := range testCase.imgs {
						if elem.ID == dres.GetId() {
							return true
						}
					}
					return false
				})

				assert.Assert(t, !exists)
			})

			t.Run("Assert Notifications", func(t *testing.T) {
				/*
					its easier to just check for length, since the order received could be different
					depending on which go routien finished first
				*/
				assert.Equal(t, len(notifChan), len(testCase.notifs))

				for range testCase.notifs {
					select {
					case notif := <-notifChan:
						found := slices.ContainsFunc(testCase.notifs, func(elem string) bool {
							return elem == notif.msg
						})

						assert.Check(t, found)
					default:
						t.Errorf("No notification received")
					}
				}
			})
		})
	}
}

// I do relise, this is not an exhaustive test. I don't understand how the delete mechanism works for volumes yet.
func TestDeleteVolume(t *testing.T) {
	tests := []struct {
		Id        string
		notifWant string
		errorStr  string
		force     bool
	}{
		{
			Id:        "1",
			notifWant: listStatusMessageStyle.Render("Deleted"),
			force:     false,
		},
	}

	mock := setupMockDockerClient(t)
	mock.ToggleContainerListAll()

	for _, testCase := range tests {
		t.Run("Force Delete Exising image", func(t *testing.T) {

			notifChan := make(chan notificationMetadata, 10)
			op := volumeDelete(mock, testCase.Id, testCase.force, 2, notifChan)

			err := op()

			// test for error
			if testCase.errorStr != "" {
				assert.ErrorContains(t, err, testCase.errorStr)
				// if there is an error, return early so that we do not perform other subtests
				return
			}

			t.Run("Confirm image deleted", func(t *testing.T) {
				images := mock.ListImages()

				exists := slices.ContainsFunc(images, func(elem it.ImageSummary) bool {
					if elem.ID == testCase.Id {
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

func TestToggleListAll(t *testing.T) {
	t.Run("Docker", func(t *testing.T) {
		dockerClient := setupMockDockerClient(t)
		assert.Assert(t, !dockerClient.GetListOptions().All)
		dockerClient.ToggleContainerListAll()
		t.Log(dockerClient.GetListOptions().All)
		assert.Assert(t, dockerClient.GetListOptions().All)
	})

	t.Run("Podman", func(t *testing.T) {
		podmanClient := setupMockDockerClient(t)
		assert.Assert(t, !podmanClient.GetListOptions().All)
		podmanClient.ToggleContainerListAll()
		t.Log(podmanClient.GetListOptions().All)
		assert.Assert(t, podmanClient.GetListOptions().All)
	})
}
