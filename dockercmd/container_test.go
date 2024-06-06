package dockercmd

import (
	"testing"
)

var dockerclient = NewDockerClient()

func BenchmarkContainerList(b *testing.B) {
	b.Run("Showing container size", func(b *testing.B) {
		for range b.N {
			dockerclient.ListContainers(false)
		}
	})
	b.Run("NOT Showing container size", func(b *testing.B) {
		for range b.N {
			dockerclient.ListContainers(true)
		}
	})
}
