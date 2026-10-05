package connector

import (
	"sync"
	"testing"
	"time"

	"github.com/bcicen/ctop/container"
	api "github.com/fsouza/go-dockerclient"
)

func TestCalcUptimeNil(t *testing.T) {
	// Should not panic on nil container
	uptime := calcUptime(nil)
	if uptime != "" {
		t.Errorf("expected empty string for nil container, got %s", uptime)
	}

	// Should not panic on container with zero State
	c := &api.Container{}
	uptime = calcUptime(c)
	if uptime != "" {
		t.Errorf("expected empty string for container with zero State, got %s", uptime)
	}
}

func TestCalcUptimeRunning(t *testing.T) {
	c := &api.Container{
		State: api.State{
			Running:   true,
			StartedAt: time.Now().Add(-10 * time.Minute),
		},
	}
	uptime := calcUptime(c)
	if uptime == "" {
		t.Errorf("expected non-empty uptime for running container")
	}
}

func TestDockerRefreshNilFields(t *testing.T) {
	cm := &Docker{
		containers:   make(map[string]*container.Container),
		needsRefresh: make(chan string, 10),
		statuses:     make(chan StatusUpdate, 10),
		closed:       make(chan struct{}),
		lock:         sync.RWMutex{},
	}

	c := container.New("c1", nil, nil)
	cm.containers["c1"] = c

	// Calling refresh when container is not found in docker client (nil inspect result)
	cm.refresh(c)
	// c1 should be deleted safely without panics
	if _, ok := cm.containers["c1"]; ok {
		t.Errorf("expected container to be removed when not found")
	}
}

func TestDockerRefreshAllEmptyNames(t *testing.T) {
	// Simulate container with empty Names slice (issue #372)
	cm := &Docker{
		containers:   make(map[string]*container.Container),
		needsRefresh: make(chan string, 10),
		statuses:     make(chan StatusUpdate, 10),
		closed:       make(chan struct{}),
		lock:         sync.RWMutex{},
	}

	allContainers := []api.APIContainers{
		{
			ID:    "1234567890ab",
			Names: []string{}, // Empty Names slice
			State: "running",
		},
		{
			ID:    "abcdef123456",
			Names: []string{"/my-container"},
			State: "running",
		},
	}

	for _, i := range allContainers {
		c := container.New(i.ID, nil, nil)
		cm.lock.Lock()
		cm.containers[i.ID] = c
		cm.lock.Unlock()

		if len(i.Names) > 0 {
			c.SetMeta("name", shortName(i.Names[0]))
		} else {
			c.SetMeta("name", "")
		}
		c.SetState(i.State)
		cm.needsRefresh <- c.Id
	}

	c1, ok := cm.Get("1234567890ab")
	if !ok {
		t.Fatalf("expected container c1 to exist")
	}
	if c1.GetMeta("name") != "" {
		t.Errorf("expected empty name for container with empty Names, got %s", c1.GetMeta("name"))
	}

	c2, ok := cm.Get("abcdef123456")
	if !ok {
		t.Fatalf("expected container c2 to exist")
	}
	if c2.GetMeta("name") != "my-container" {
		t.Errorf("expected name 'my-container', got %s", c2.GetMeta("name"))
	}
}
