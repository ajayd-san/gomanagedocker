package dockercmd

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListContainer(t *testing.T) {

	client := NewDockerClient()

	containersList := client.ListContainers()
	indenttest, _ := json.MarshalIndent(containersList, "", "\t")
	log.Println(string(indenttest))
	assert.Equal(t, len(containersList), 1, "Not all containers detected")
}
