package dockercmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListContainer(t *testing.T) {

	client := NewDockerClient()

	containersList := client.ListContainers()
	assert.Equal(t, len(containersList), 1, "Not all containers detected")
}
