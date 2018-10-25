package collector

import (
	"bufio"
	"context"
	"io"
	"strings"
	"time"

	"github.com/bcicen/ctop/models"
	api "github.com/fsouza/go-dockerclient"
)

type DockerLogs struct {
	id     string
	client *api.Client
	done   chan bool
}

func NewDockerLogs(id string, client *api.Client) *DockerLogs {
	return &DockerLogs{
		id:     id,
		client: client,
		done:   make(chan bool),
	}
}

func (l *DockerLogs) Stream() chan models.Log {
	r, w := io.Pipe()
	logCh := make(chan models.Log)
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
			parts := strings.Split(scanner.Text(), " ")
			ts := l.parseTime(parts[0])
			logCh <- models.Log{ts, strings.Join(parts[1:], " ")}
		}
	}()

	// connect to container log stream
	go func() {
		err := l.client.Logs(opts)
		if err != nil {
			log.Errorf("error reading container logs: %s", err)
		}
		log.Infof("log reader stopped for container: %s", l.id)
	}()

	go func() {
		<-l.done
		cancel()
	}()

	log.Infof("log reader started for container: %s", l.id)
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
