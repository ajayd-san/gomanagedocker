package dockercmd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseDockerScoutOutput(t *testing.T) {
	type test struct {
		input string
		want  ScoutData
	}

	cases := []test{
		{
			`
          Target             │  nginx:latest          │    0C     0H     1M    48L     1?
            digest           │  1445eb9c6dc5          │
          Base image         │  debian:bookworm-slim  │    0C     0H     0M    23L
          Updated base image │  debian:stable-slim    │    0C     0H     0M    23L
                             │                        │

        What's next:
            Include policy results in your quickview by supplying an organization → docker scout quickview nginx --org <organization>
	`,
			ScoutData{
				[]ImageVulnerabilities{
					{Label: "Target", ImageName: "nginx:latest", Critical: "0", High: "0", Medium: "1", Low: "48", UnknownSeverity: "1"},
					{Label: "Base image", ImageName: "debian:bookworm-slim", Critical: "0", High: "0", Medium: "0", Low: "23", UnknownSeverity: ""},
					{Label: "Updated base image", ImageName: "debian:stable-slim", Critical: "0", High: "0", Medium: "0", Low: "23", UnknownSeverity: ""},
				},
			},
		},
		{
			`

          Target             │  myimage:latest          │    1C     2H     3M    4L     5?
            digest           │  abcdef123456			│
          Base image         │  ubuntu:20.04			│    1C     2H     3M    4L
          Updated base image │  ubuntu:latest			│    0C     0H     1M    2L
                             │							│
`,
			ScoutData{
				[]ImageVulnerabilities{
					{"Target", "myimage:latest", "1", "2", "3", "4", "5"},
					{"Base image", "ubuntu:20.04", "1", "2", "3", "4", ""},
					{"Updated base image", "ubuntu:latest", "0", "0", "1", "2", ""},
				},
			},
		},
	}

	for _, tcase := range cases {
		got := parseDockerScoutOutput([]byte(tcase.input))

		if !cmp.Equal(&tcase.want, got) {
			t.Fatalf("structs do not match\n %s", cmp.Diff(&tcase.want, got))
		}
	}

}
