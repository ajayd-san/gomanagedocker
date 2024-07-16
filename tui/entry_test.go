package tui

import (
	"testing"

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
