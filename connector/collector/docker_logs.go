package collector

import (
	"bufio"
	"context"
	"io"

	api "github.com/fsouza/go-dockerclient"
)

type DockerLogs struct {
	id     string
	client *api.Client
	done   chan bool
}

func (l *DockerLogs) Stream() chan string {
	r, w := io.Pipe()
	logCh := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())

	opts := api.LogsOptions{
		Context:      ctx,
		Container:    l.id,
		OutputStream: w,
		ErrorStream:  w,
		Stdout:       true,
		Stderr:       true,
		Tail:         "10",
		Follow:       true,
		Timestamps:   true,
	}

	// read io pipe into channel
	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			logCh <- scanner.Text()
		}
	}()

	// connect to container log stream
	go func() {
		err := l.client.Logs(opts)
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
