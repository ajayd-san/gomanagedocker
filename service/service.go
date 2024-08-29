package service

import (
	"os/exec"

	"github.com/ajayd-san/gomanagedocker/service/types"
	et "github.com/docker/docker/api/types"
)

// both DockerClient and PodmanClient satisfy this interface
type Service interface {
	Ping() error
	GetListOptions() types.ContainerListOptions

	// image
	BuildImage(buildContext string, options et.ImageBuildOptions) (*et.ImageBuildResponse, error)
	ListImages() []types.ImageSummary
	RunImage(config types.ContainerCreateConfig) (*string, error)
	DeleteImage(id string, opts types.RemoveImageOptions) error
	PruneImages() (types.ImagePruneReport, error)

	// container
	InspectContainer(id string) (*types.InspectContainerData, error)
	ListContainers(showContainerSize bool) []types.ContainerSummary
	ToggleContainerListAll()
	ToggleStartStopContainer(id string, isRunning bool) error
	RestartContainer(id string) error
	TogglePauseResume(id string, state string) error
	DeleteContainer(id string, opts types.ContainerRemoveOpts) error
	PruneContainers() (types.ContainerPruneReport, error)
	ExecCmd(id string) *exec.Cmd
	LogsCmd(id string) *exec.Cmd

	// volume
	ListVolumes() ([]types.VolumeSummary, error)
	PruneVolumes() (*types.VolumePruneReport, error)
	DeleteVolume(id string, force bool) error
}
