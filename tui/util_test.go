package tui

import (
	"strings"
	"testing"

	"github.com/ajayd-san/gomanagedocker/service/dockercmd"
	it "github.com/ajayd-san/gomanagedocker/service/types"
	"github.com/docker/docker/api/types"
	"gotest.tools/v3/assert"
)

func TestNotifyList(t *testing.T) {
	api := dockercmd.MockApi{}

	containers := []types.Container{
		{
			Names:      []string{"a"},
			ID:         "1",
			SizeRw:     1e+9,
			SizeRootFs: 2e+9,
			State:      "running",
			Status:     "",
		},
	}

	api.SetMockContainers(containers)

	mockcli := dockercmd.NewMockCli(&api)

	keymap := NewKeyMap(it.Docker)

	CONTAINERS = 0
	model := MainModel{
		dockerClient: mockcli,
		activeTab:    0,
		TabContent: []listModel{
			InitList(0, keymap.container, keymap.containerBulk),
		},
	}

	t.Run("Notify test", func(t *testing.T) {
		NotifyList(model.getActiveList(), "Kiryu")
		got := model.View()
		contains := "Kiryu"
		assert.Check(t, strings.Contains(got, contains))
	})
}

func TestSepPortMapping(t *testing.T) {
	t.Run("Clean string, test mapping", func(t *testing.T) {
		// format is host:container
		testStr := "8080:80/tcp,1123:112,6969:9696/udp"
		want := []it.PortBinding{
			{
				HostPort:      "8080",
				ContainerPort: "80",
				Proto:         "tcp",
			},
			{
				HostPort:      "1123",
				ContainerPort: "112",
				Proto:         "tcp",
			},
			{
				HostPort:      "6969",
				ContainerPort: "9696",
				Proto:         "udp",
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
