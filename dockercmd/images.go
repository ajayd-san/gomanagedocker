package dockercmd

import (
	"context"
	"log"

	"github.com/docker/docker/api/types/image"
)

func (dc *DockerClient) ListImages() []image.Summary {
	images, err := dc.cli.ImageList(context.Background(), image.ListOptions{})

	if err != nil {
		panic(err)
	}

	return images
}

func (dc *DockerClient) DeleteImage(id string) error {
	res, err := dc.cli.ImageRemove(context.Background(), id, image.RemoveOptions{})
	log.Println(res)
	return err
}
