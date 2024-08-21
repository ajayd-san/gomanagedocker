package service

import (
	"github.com/ajayd-san/gomanagedocker/service/types"
	et "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

// both DockerClient and PodmanClient satisfy this interface
type Service interface {
	Ping() error
	GetListOptions() *container.ListOptions

	// image
	BuildImage(buildContext string, options et.ImageBuildOptions) (*et.ImageBuildResponse, error)
	ListImages() []types.ImageSummary
	RunImage(containerConfig *container.Config, hostConfig *container.HostConfig, containerName string) (*string, error)
	DeleteImage(id string, opts types.RemoveImageOptions) error
	PruneImages() (et.ImagesPruneReport, error)

	// container
	InspectContainer(id string) (*types.InspectContainerData, error)
	ListContainers(showContainerSize bool) []types.ContainerSummary
	ToggleContainerListAll()
	ToggleStartStopContainer(id string) error
	RestartContainer(id string) error
	TogglePauseResume(id string) error
	DeleteContainer(id string, opts container.RemoveOptions) error
	PruneContainers() (et.ContainersPruneReport, error)

	// volume
	ListVolumes() ([]types.VolumeSummary, error)
	PruneVolumes() (*et.VolumesPruneReport, error)
	DeleteVolume(id string, force bool) error
}
