package tui

import (
	"strings"
	"testing"

	"github.com/ajayd-san/gomanagedocker/dockercmd"
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

	CONTAINERS = 0
	model := MainModel{
		dockerClient: mockcli,
		activeTab:    0,
		TabContent: []listModel{
			InitList(0),
		},
	}

	t.Run("Notify test", func(t *testing.T) {
		NotifyList(model.getActiveList(), "Kiryu")
		got := model.View()
		contains := "Kiryu"
		assert.Check(t, strings.Contains(got, contains))
	})
}
