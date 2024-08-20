package dockercmd

import (
	"bytes"
	"testing"

	"gotest.tools/v3/assert"
)

func TestGetDockerIgnorePatternsFromFS(t *testing.T) {
	t.Run("Test 1", func(t *testing.T) {
		file := `abc
czy*
*/bar
`
		content := bytes.NewBuffer([]byte(file))
		got := getDockerIgnorePatterns(content)

		want := []string{"abc", "czy*", "*/bar"}
		assert.DeepEqual(t, got, want)

	})

}
