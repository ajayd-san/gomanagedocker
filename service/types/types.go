package types

import "io"

type ServiceType int

const (
	Docker ServiceType = iota
	Podman
)

type ImageSummary struct {
	ID         string
	Size       int64
	RepoTags   []string
	Containers int64
	Created    int64
}

/*
this type direct copy of podman's `types.RemoveImageOptions`,
I chose this cuz it is more exhausive compared to docker's
*/
type RemoveImageOptions struct {
	All            bool
	Force          bool
	Ignore         bool
	LookupManifest bool
	NoPrune        bool
}

type ImagePruneReport struct {
	ImagesDeleted int
}

type ImageBuildOptions struct {
	Tags       []string
	Dockerfile string
}

type ImageBuildReport struct {
	Body io.Reader
}

type ImageBuildJSON struct {
	Stream string     `json:"stream,omitempty"`
	Error  *JSONError `json:"errorDetail,omitempty"`
}

type JSONError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type ContainerSummary struct {
	ID      string
	ImageID string
	Created int64
	Names   []string
	State   string
	Command string
	// Status     string
	SizeRw     int64
	SizeRootFs int64
	// Mounts []string
	// Ports []string
}

// // represents container state
// type ContainerState struct {
// 	Status     string // String representation of the container state. Can be one of "created", "running", "paused", "restarting", "removing", "exited", or "dead"
// 	Running    bool
// 	Paused     bool
// 	Restarting bool
// 	OOMKilled  bool
// 	Dead       bool
// 	Pid        int
// 	ExitCode   int
// 	Error      string
// }

type VolumeSummary struct {
	Name       string
	CreatedAt  string
	Driver     string
	Mountpoint string
	UsageData  int64
}

type VolumePruneReport struct {
	VolumesPruned int
}

type InspectContainerData struct {
	ContainerSummary
}

type ContainerListOptions struct {
	All  bool
	Size bool
}

type ContainerRemoveOpts struct {
	Force         bool
	RemoveVolumes bool
	RemoveLinks   bool
}

type ContainerPruneReport struct {
	ContainersDeleted int
}

type ContainerCreateConfig struct {
	// ExposedPorts []PortMapping
	// name of the container
	Name string
	Env  []string
	// ID of image
	ImageId      string
	PortBindings []PortBinding
}

type PortBinding struct {
	HostPort      string
	ContainerPort string
	Proto         string
}
