package dockercmd

import (
	"context"

	"github.com/docker/docker/api/types/image"
)

func (dc *DockerClient) ListImages() []image.Summary {
	images, err := dc.cli.ImageList(context.Background(), image.ListOptions{})

	if err != nil {
		panic(err)
	}

	return images
}
