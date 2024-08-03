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
	//this makes sure "" is not printed in the table
	unknownSev := string(submatches[8])
	if unknownSev == "" {
		unknownSev = "0"
	}

	return ImageVulnerabilities{
		Label:           string(submatches[1]),
		ImageName:       string(submatches[2]),
		Critical:        string(submatches[3]),
		High:            string(submatches[4]),
		Medium:          string(submatches[5]),
		Low:             string(submatches[6]),
		UnknownSeverity: unknownSev,
	}

}

type ScoutData struct {
	ImageVulEntries []ImageVulnerabilities
}

type DockerClient struct {
	cli               client.CommonAPIClient
	containerListArgs container.ListOptions
}

func NewDockerClient() DockerClient {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	//TODO: size should not be true, investigate later
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

// used for testing only
func NewMockCli(cli *MockApi) DockerClient {
	return DockerClient{
		cli:               cli,
		containerListArgs: container.ListOptions{},
	}
}

// util
func (dc DockerClient) GetListOptions() *container.ListOptions {
	return &dc.containerListArgs
}
