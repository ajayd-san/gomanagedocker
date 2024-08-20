package dockercmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

func timeBenchmark(start time.Time, msg string) {
	timeTook := time.Since(start)
	log.Println(fmt.Sprintf("%s : %s", msg, timeTook))
}

func getDockerIgnorePatterns(file io.Reader) []string {
	patterns := make([]string, 0)
	buffer := bufio.NewReader(file)

	for {
		line, err := buffer.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSuffix(line, "\n")
		patterns = append(patterns, line)
	}

	return patterns
}
