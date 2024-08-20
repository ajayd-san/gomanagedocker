package service

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/volume"
)

// both DockerClient and PodmanClient satisfy this interface
type Service interface {
	// image
	BuildImage(buildContext string, options types.ImageBuildOptions) (*types.ImageBuildResponse, error)
	ListImages() []image.Summary
	RunImage(containerConfig *container.Config, hostConfig *container.HostConfig, containerName string) (*string, error)
	DeleteImage(id string, opts image.RemoveOptions) error
	PruneImages() (types.ImagesPruneReport, error)

	// container
	ListContainers(showContainerSize bool) []types.Container
	ToggleContainerListAll()
	ToggleStartStopContainer(id string) error
	RestartContainer(id string) error
	TogglePauseResume(id string) error
	DeleteContainer(id string, opts container.RemoveOptions) error
	PruneContainers() (types.ContainersPruneReport, error)

	// volume
	ListVolumes() ([]*volume.Volume, error)
	PruneVolumes() (*types.VolumesPruneReport, error)
	DeleteVolume(id string, force bool) error
}
