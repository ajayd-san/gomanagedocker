package podmancmd

import (
	"strconv"

	it "github.com/ajayd-san/gomanagedocker/service/types"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/bindings/images"
	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/docker/docker/api/types"

	nettypes "github.com/containers/common/libnetwork/types"
)

func (po *PodmanClient) BuildImage(buildContext string, options types.ImageBuildOptions) (*types.ImageBuildResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (pc *PodmanClient) ListImages() []it.ImageSummary {
	raw, err := images.List(pc.cli, nil)

	if err != nil {
		panic(err)
	}

	return toImageSummaryArr(raw)
}

// runs image and returns container ID
func (pc *PodmanClient) RunImage(config it.ContainerCreateConfig) (*string, error) {
	spec := specgen.NewSpecGenerator(config.ImageId, false)

	envMap, err := getEnvMap(&config.Env)

	if err != nil {
		return nil, err
	}

	bindings := make([]nettypes.PortMapping, len(config.PortBindings))
	for i, mapping := range config.PortBindings {
		containerPort, _ := strconv.ParseUint(mapping.ContainerPort, 10, 16)
		HostPort, _ := strconv.ParseUint(mapping.HostPort, 10, 16)

		bindings[i] = nettypes.PortMapping{
			HostIP:        "::1",
			ContainerPort: uint16(containerPort),
			HostPort:      uint16(HostPort),
			Protocol:      mapping.Proto,
		}
	}

	spec.Name = config.Name
	spec.Env = envMap
	spec.PortMappings = bindings
	spec.NetNS = specgen.Namespace{
		NSMode: specgen.Bridge,
	}

	res, err := containers.CreateWithSpec(pc.cli, spec, nil)

	if err != nil {
		return nil, err
	}

	err = containers.Start(pc.cli, res.ID, nil)

	if err != nil {
		return nil, err
	}

	return &res.ID, nil
}

func (pc *PodmanClient) DeleteImage(id string, opts it.RemoveImageOptions) error {
	_, errs := images.Remove(pc.cli, []string{id}, &images.RemoveOptions{
		All:            &opts.All,
		Force:          &opts.Force,
		Ignore:         &opts.Ignore,
		LookupManifest: &opts.LookupManifest,
		NoPrune:        &opts.NoPrune,
	})

	if errs != nil {
		return errs[0]
	}

	return nil
}

func (pc *PodmanClient) PruneImages() (it.ImagePruneReport, error) {
	t := true
	reports, err := images.Prune(pc.cli, &images.PruneOptions{
		All: &t,
	})

	return it.ImagePruneReport{ImagesDeleted: len(reports)}, err
}
