package dockercmd

import (
	"bufio"
	"os"
	"slices"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dimage "github.com/docker/docker/api/types/image"
	"github.com/google/go-cmp/cmp"
	"gotest.tools/v3/assert"
)

func TestListImages(t *testing.T) {

	imgs := []dimage.Summary{
		{
			Containers: 0,
			ID:         "0",
		},

		{
			Containers: 2,
			ID:         "1",
		},
		{
			Containers: 3,
			ID:         "2",
		},
		{
			Containers: 4,
			ID:         "5",
		},
	}

	dclient := DockerClient{
		cli: &MockApi{
			mockImages:      imgs,
			CommonAPIClient: nil,
		},
		containerListArgs: container.ListOptions{},
	}

	got := dclient.ListImages()
	want := imgs

	assert.DeepEqual(t, got, want)
}

func TestDeleteImage(t *testing.T) {
	imgs := []dimage.Summary{
		{
			Containers: 0,
			ID:         "0",
		},

		{
			Containers: 2,
			ID:         "1",
		},
		{
			Containers: 3,
			ID:         "2",
		},
		{
			Containers: 4,
			ID:         "5",
		},
	}

	dclient := DockerClient{
		cli: &MockApi{
			mockImages:      imgs,
			CommonAPIClient: nil,
		},
		containerListArgs: container.ListOptions{},
	}

	t.Run("No force required image test", func(t *testing.T) {
		err := dclient.DeleteImage("0", dimage.RemoveOptions{})
		assert.NilError(t, err)

		afterDeleteImgs := dclient.cli.(*MockApi).mockImages

		// we do len(img) - 1 because the slices.Delete swaps the 'to be deleted' index with last index and zeros it. so we exclude the last element in the array
		assert.DeepEqual(t, afterDeleteImgs, imgs[0:len(imgs)-1])
	})

	t.Run("Should fail, image has active containers", func(t *testing.T) {
		err := dclient.DeleteImage("1", dimage.RemoveOptions{})
		assert.ErrorContains(t, err, "must be forced")
	})

	t.Run("With force", func(t *testing.T) {
		err := dclient.DeleteImage("1", dimage.RemoveOptions{Force: true})
		assert.NilError(t, err)

		// same reason as above, but this time we exclude last to elements
		afterDeleteImgs := dclient.cli.(*MockApi).mockImages
		assert.DeepEqual(t, afterDeleteImgs, imgs[0:len(imgs)-2])
	})
}

func TestPruneImages(t *testing.T) {
	imgs := []dimage.Summary{
		{
			Containers: 0,
			ID:         "0",
		},

		{
			Containers: 0,
			ID:         "1",
		},
		{
			Containers: 3,
			ID:         "2",
		},
		{
			Containers: 0,
			ID:         "5",
		},
	}

	dclient := DockerClient{
		cli: &MockApi{
			mockImages:      imgs,
			CommonAPIClient: nil,
		},
		containerListArgs: container.ListOptions{},
	}

	dclient.PruneImages()

	finalImages := dclient.cli.(*MockApi).mockImages
	want := []dimage.Summary{
		{
			Containers: 3,
			ID:         "2",
		},
	}
	assert.DeepEqual(t, finalImages, want)
}
func TestParseDockerScoutOutput(t *testing.T) {
	type test struct {
		input string
		want  ScoutData
	}

	cases := []test{
		{
			`
          Target             │  nginx:latest          │    0C     0H     1M    48L     1?
            digest           │  1445eb9c6dc5          │
          Base image         │  debian:bookworm-slim  │    0C     0H     0M    23L
          Updated base image │  debian:stable-slim    │    0C     0H     0M    23L
                             │                        │

        What's next:
            Include policy results in your quickview by supplying an organization → docker scout quickview nginx --org <organization>
	`,
			ScoutData{
				[]ImageVulnerabilities{
					{Label: "Target", ImageName: "nginx:latest", Critical: "0", High: "0", Medium: "1", Low: "48", UnknownSeverity: "1"},
					{Label: "Base image", ImageName: "debian:bookworm-slim", Critical: "0", High: "0", Medium: "0", Low: "23", UnknownSeverity: "0"},
					{Label: "Updated base image", ImageName: "debian:stable-slim", Critical: "0", High: "0", Medium: "0", Low: "23", UnknownSeverity: "0"},
				},
			},
		},
		{
			`

          Target             │  myimage:latest          │    1C     2H     3M    4L     5?
            digest           │  abcdef123456			│
          Base image         │  ubuntu:20.04			│    1C     2H     3M    4L
          Updated base image │  ubuntu:latest			│    0C     0H     1M    2L
                             │							│
`,
			ScoutData{
				[]ImageVulnerabilities{
					{"Target", "myimage:latest", "1", "2", "3", "4", "5"},
					{"Base image", "ubuntu:20.04", "1", "2", "3", "4", "0"},
					{"Updated base image", "ubuntu:latest", "0", "0", "1", "2", "0"},
				},
			},
		},
	}

	for _, tcase := range cases {
		got := parseDockerScoutOutput([]byte(tcase.input))

		if !cmp.Equal(&tcase.want, got) {
			t.Fatalf("structs do not match\n %s", cmp.Diff(&tcase.want, got))
		}
	}

}

func TestBuildImage(t *testing.T) {
	imgs := []dimage.Summary{
		{
			Containers: 0,
			ID:         "0",
		},

		{
			Containers: 0,
			ID:         "1",
		},
		{
			Containers: 3,
			ID:         "2",
		},
		{
			Containers: 0,
			ID:         "5",
		},
	}

	dclient := DockerClient{
		cli: &MockApi{
			mockImages:      imgs,
			CommonAPIClient: nil,
		},
		containerListArgs: container.ListOptions{},
	}

	cwd, _ := os.Getwd()
	opts := types.ImageBuildOptions{
		Tags: []string{"test"},
	}
	res, err := dclient.BuildImage(cwd, opts)
	if err != nil {
		t.Error(err)
	}

	// no-op, must wait till this finishes
	reader := bufio.NewScanner(res.Body)
	for reader.Scan() {
	}

	got := dclient.ListImages()

	index := slices.IndexFunc(got, func(entry dimage.Summary) bool {
		return slices.Equal(entry.RepoTags, []string{"test"})
	})

	if index == -1 {
		t.Error("Could not find built image")
	}
}

func TestSepPortMapping(t *testing.T) {
	t.Run("Clean string, test mapping", func(t *testing.T) {
		// format is host:container
		testStr := "8080:80/tcp,1123:112,6969:9696/udp"
		want := []PortBinding{
			{
				HostPort:      "8080",
				ContainerPort: "80",
				Proto:         "tcp",
			},
			{
				"1123",
				"112",
				"tcp",
			},
			{
				"6969",
				"9696",
				"udp",
			},
		}

		got, err := GetPortMappingFromStr(testStr)

		assert.NilError(t, err)

		assert.DeepEqual(t, got, want)
	})

	t.Run("Empty port string", func(t *testing.T) {
		testStr := ""
		_, err := GetPortMappingFromStr(testStr)
		assert.NilError(t, err)
	})

	t.Run("Invalid mapping, should throw error", func(t *testing.T) {
		testStr := "8080:878:9/tcp"
		_, err := GetPortMappingFromStr(testStr)
		assert.Error(t, err, "Port Mapping 8080:878:9/tcp is invalid")
	})
}
