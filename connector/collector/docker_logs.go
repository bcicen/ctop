package collector

import (
	"bufio"
	"context"
	"io"
	"strings"
	"time"

	"github.com/bcicen/ctop/models"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
)

type DockerLogs struct {
	id     string
	client *client.Client
	done   chan bool
}

func (l *DockerLogs) Stream() chan models.Log {
	r, _ := io.Pipe()
	logCh := make(chan models.Log)
	ctx, cancel := context.WithCancel(context.Background())

	opts := types.ContainerLogsOptions{
		ShowStdout:   true,
		ShowStderr:   true,
		Tail:         "10",
		Follow:       true,
		Timestamps:   true,
	}

	// read io pipe into channel
	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			parts := strings.Split(scanner.Text(), " ")
			ts := l.parseTime(parts[0])
			logCh <- models.Log{ts, strings.Join(parts[1:], " ")}
		}
	}()

	// connect to container log stream
	go func() {
		_, err := l.client.ContainerLogs(ctx, l.id, opts)
		if err != nil {
			log.Errorf("error reading container logs: %s", err)
		}
	}()

	go func() {
		select {
		case <-l.done:
			cancel()
		}
	}()

	return logCh
}

func (l *DockerLogs) Stop() { l.done <- true }

func (l *DockerLogs) parseTime(s string) time.Time {
	ts, err := time.Parse("2006-01-02T15:04:05.000000000Z", s)
	if err != nil {
		log.Errorf("failed to parse container log: %s", err)
		ts = time.Now()
	}
	return ts
}
