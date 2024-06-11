package dockercmd

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type ImageVulnerabilities struct {
	Label           string
	ImageName       string
	Critical        string
	High            string
	Medium          string
	Low             string
	UnknownSeverity string
}

// takes [][]bytes returned by regex.FindSubmatches and returns ImageVulnerabilities
func makeImageVulnerabilities(submatches [][]byte) ImageVulnerabilities {
	return ImageVulnerabilities{
		Label:           string(submatches[1]),
		ImageName:       string(submatches[2]),
		Critical:        string(submatches[3]),
		High:            string(submatches[4]),
		Medium:          string(submatches[5]),
		Low:             string(submatches[6]),
		UnknownSeverity: string(submatches[8]),
	}

}

type ScoutData struct {
	ImageVulEntries []ImageVulnerabilities
	TargetDigest    string
}

type DockerClient struct {
	cli               *client.Client
	containerListArgs container.ListOptions
}

func NewDockerClient() DockerClient {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return DockerClient{
		cli: cli,
		containerListArgs: container.ListOptions{
			Size:   true,
			All:    false,
			Latest: false,
		},
	}
}

func (dc DockerClient) PingDocker() error {
	_, err := dc.cli.Ping(context.Background())
	return err
}
