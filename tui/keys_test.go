package tui

import (
	"testing"

	it "github.com/ajayd-san/gomanagedocker/service/types"
	"gotest.tools/v3/assert"
)

func TestNewKeyMap(t *testing.T) {
	t.Run("Docker, Scout should be enabled enabled", func(t *testing.T) {
		dockerKeymap := NewKeyMap(it.Docker)
		scoutEnabled := dockerKeymap.image.Scout.Enabled()
		assert.Assert(t, scoutEnabled)
	})

	t.Run("Podman, Scout should be disabled", func(t *testing.T) {
		dockerKeymap := NewKeyMap(it.Podman)
		scoutEnabled := dockerKeymap.image.Scout.Enabled()
		assert.Assert(t, !scoutEnabled)
	})

}
