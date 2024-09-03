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

func TestCustomWriter(t *testing.T) {
	sub := CustomWriter{}

	data := `{
		"name": "idk",
		"age": "fucker"
	}`

	i, err := sub.Write([]byte(data))
	assert.NilError(t, err)

	t.Log(i)

	got := make([]byte, i)
	_, err = sub.Read(got)
	t.Log(got)

	assert.NilError(t, err)
	assert.DeepEqual(t, got, []byte(data))
}
