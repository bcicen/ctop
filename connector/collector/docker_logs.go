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
		//ErrorStream:  w,
		Stdout:      true,
		Stderr:      true,
		Tail:        "20",
		Follow:      true,
		Timestamps:  true,
		RawTerminal: true,
	}

	// read io pipe into channel
	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			parts := strings.SplitN(scanner.Text(), " ", 2)
			if len(parts) == 0 {
				continue
			}
			if len(parts) < 2 {
				logCh <- models.Log{Timestamp: l.parseTime(""), Message: parts[0]}
			} else {
				logCh <- models.Log{Timestamp: l.parseTime(parts[0]), Message: parts[1]}
			}
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
	ts, err := time.Parse(time.RFC3339Nano, s)
	if err == nil {
		return ts
	}

	ts, err2 := time.Parse(time.RFC3339Nano, l.stripPfx(s))
	if err2 == nil {
		return ts
	}

	log.Errorf("failed to parse container log: %s", err)
	log.Errorf("failed to parse container log2: %s", err2)
	return time.Now()
}

// attempt to strip message header prefix from a given raw docker log string
func (l *DockerLogs) stripPfx(s string) string {
	b := []byte(s)
	if len(b) > 8 {
		return string(b[8:])
	}
	return s
}
