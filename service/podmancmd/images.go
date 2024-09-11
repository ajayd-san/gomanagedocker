package podmancmd

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"strconv"

	it "github.com/ajayd-san/gomanagedocker/service/types"
	"github.com/containers/buildah/define"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/bindings/images"
	"github.com/containers/podman/v5/pkg/domain/entities/types"
	"github.com/containers/podman/v5/pkg/specgen"

	nettypes "github.com/containers/common/libnetwork/types"
)

func (pc *PodmanClient) BuildImage(buildContext string, options it.ImageBuildOptions) (*it.ImageBuildReport, error) {

	/*
		INFO: this method has a lot going on, we return a io.Reader that receives data in form of types.ImageBuildJSON.
			We want to do this as the parallelly as each step in dockerfile gets processed by image.Build which is why
			we use pipes.
	*/
	outR, outW := io.Pipe()
	reportPipeR, reportPipeW := io.Pipe()

	reportReader := bufio.NewReader(reportPipeR)

	// we use this to send the error from the builder goroutine to the reader goroutine(below)
	errChan := make(chan error, 2)

	go func() {
		reader := bufio.NewReader(outR)
		for {
			str, err := reader.ReadString('\n')
			if err != nil {
				break
			}

			var bytes []byte

			step := it.ImageBuildJSON{
				Stream: str,
			}
			bytes, _ = json.Marshal(step)

			if err != nil {
				log.Printf("Marshalling Error: %s", err.Error())
			}
			reportPipeW.Write(bytes)
		}

		select {
		case err := <-errChan:
			if err != nil {
				errReport := it.ImageBuildJSON{
					Error: &it.JSONError{
						Message: err.Error(),
					},
				}

				bytes, _ := json.Marshal(errReport)
				reportPipeW.Write(bytes)
			}
		}

		reportPipeW.Close()
	}()

	go func() {
		//TODO: registry option
		_, err := images.Build(pc.cli, []string{options.Dockerfile}, types.BuildOptions{
			BuildOptions: define.BuildOptions{
				// Labels:         []string{"teststr"},
				// Registry:       "regname",
				AdditionalTags: options.Tags,
				Out:            outW,
			},
		})

		errChan <- err

		outW.Close()

	}()

	return &it.ImageBuildReport{
		Body: reportReader,
	}, nil

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
