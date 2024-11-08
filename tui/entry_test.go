package tui

import (
	"testing"

	"github.com/ajayd-san/gomanagedocker/service/types"
	"gotest.tools/v3/assert"
)

func TestSetTabConstants(t *testing.T) {

	t.Run("Order 1", func(t *testing.T) {
		order := []string{"containers", "volumes", "images"}
		setTabConstants(order)
		assert.Equal(t, CONTAINERS, tabId(0))
		assert.Equal(t, VOLUMES, tabId(1))
		assert.Equal(t, IMAGES, tabId(2))
	})

	t.Run("Order 2", func(t *testing.T) {
		order := []string{"images", "volumes"}
		setTabConstants(order)
		assert.Equal(t, CONTAINERS, tabId(999))
		assert.Equal(t, IMAGES, tabId(0))
		assert.Equal(t, VOLUMES, tabId(1))
	})
}

func TestLoadConfig(t *testing.T) {
	t.Run("with Docker", func(t *testing.T) {
		readConfig()
		loadConfig(types.Docker)
		assert.DeepEqual(t, CONFIG_TAB_ORDERING, []string{"images", "containers", "volumes"})
	})

	t.Run("with Podman", func(t *testing.T) {
		readConfig()
		loadConfig(types.Podman)
		assert.DeepEqual(t, CONFIG_TAB_ORDERING, []string{"images", "containers", "volumes", "pods"})
	})
}
