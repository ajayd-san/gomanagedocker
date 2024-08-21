package types

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

type ContainerSummary struct {
	ID         string
	ImageID    string
	Created    int64
	Names      []string
	State      string
	Command    string
	Status     string
	SizeRw     int64
	SizeRootFs int64
	// Mounts []string
	// Ports []string
}

type InspectContainerData struct {
	ContainerSummary
}
