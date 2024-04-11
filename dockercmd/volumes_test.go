package dockercmd

import (
	"testing"
)

func TestListVolumes(t *testing.T) {

	client := NewDockerClient()

	containersList, _ := client.ListVolumes()
	for i := range containersList {
		t.Logf("%#v", containersList[i])

	}
	// assert.Equal(t, len(containersList), 1, "Not all containers detected")
}
