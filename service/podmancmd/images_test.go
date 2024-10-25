// DELETE THIS
package podmancmd

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestListImages(t *testing.T) {
	pc, err := NewPodmanClient()

	assert.NilError(t, err)

	summary := pc.ListImages()

	assert.NilError(t, err)

	assert.Assert(t, len(summary) != 0)
}
