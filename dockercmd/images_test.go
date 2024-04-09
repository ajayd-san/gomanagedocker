package dockercmd

import "testing"

func TestListImages(t *testing.T) {
	cli := NewDockerClient()

	images := cli.ListImages()

	for _, img := range images {
		t.Logf("%#v\n", img)
	}

}
