package main

import (
	"sort"
	"sync"
	"time"

	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/metrics"
	"github.com/fsouza/go-dockerclient"
)

var lock = sync.RWMutex{}

type ContainerSource interface {
	All() []*Container
	Get(string) (*Container, bool)
}

type DockerContainerSource struct {
	client       *docker.Client
	containers   Containers
	collectors   map[string]metrics.Collector
	needsRefresh map[string]int // container IDs requiring refresh
}

func NewDockerContainerSource() *DockerContainerSource {
	// init docker client
	client, err := docker.NewClient(config.GetVal("dockerHost"))
	if err != nil {
		panic(err)
	}
	cm := &DockerContainerSource{
		client:       client,
		collectors:   make(map[string]metrics.Collector),
		needsRefresh: make(map[string]int),
	}
	cm.refreshAll()
	go cm.Loop()
	go cm.watchEvents()
	return cm
}

// Docker events watcher
func (cm *DockerContainerSource) watchEvents() {
	log.Info("docker event listener starting")
	events := make(chan *docker.APIEvents)
	cm.client.AddEventListener(events)

	for e := range events {
		if e.Type != "container" {
			continue
		}
		switch e.Action {
		case "start", "die", "pause", "unpause":
			log.Debugf("handling docker event: action=%s id=%s", e.Action, e.ID)
			cm.needsRefresh[e.ID] = 1
		case "destroy":
			log.Debugf("handling docker event: action=%s id=%s", e.Action, e.ID)
			cm.delByID(e.ID)
		}
	}
}

func (cm *DockerContainerSource) refresh(id string) {
	insp := cm.inspect(id)
	// remove container if no longer exists
	if insp == nil {
		cm.delByID(id)
		return
	}

	c, ok := cm.Get(id)
	// append container struct for new containers
	if !ok {
		c = NewContainer(id, insp.Name)
		lock.Lock()
		cm.containers = append(cm.containers, c)
		lock.Unlock()
		// create collector
		if _, ok := cm.collectors[id]; ok == false {
			cm.collectors[id] = metrics.NewDocker(cm.client, id)
		}
	}

	c.SetState(insp.State.Status)

	// start collector if needed
	if c.state == "running" && !cm.collectors[c.id].Running() {
		cm.collectors[c.id].Start()
		c.Read(cm.collectors[c.id].Stream())
	}
	// stop collector if needed
	if c.state != "running" && cm.collectors[c.id].Running() {
		cm.collectors[c.id].Stop()
	}
}

func (cm *DockerContainerSource) inspect(id string) *docker.Container {
	c, err := cm.client.InspectContainer(id)
	if err != nil {
		if _, ok := err.(*docker.NoSuchContainer); ok == false {
			log.Errorf(err.Error())
		}
	}
	return c
}

// Mark all container IDs for refresh
func (cm *DockerContainerSource) refreshAll() {
	opts := docker.ListContainersOptions{All: true}
	allContainers, err := cm.client.ListContainers(opts)
	if err != nil {
		panic(err)
	}

	for _, c := range allContainers {
		cm.needsRefresh[c.ID] = 1
	}
}

func (cm *DockerContainerSource) Loop() {
	for {
		switch {
		case len(cm.needsRefresh) > 0:
			processed := []string{}
			for id, _ := range cm.needsRefresh {
				cm.refresh(id)
				processed = append(processed, id)
			}
			for _, id := range processed {
				delete(cm.needsRefresh, id)
			}
		default:
			time.Sleep(3 * time.Second)
		}
	}
}

// Get a single container, by ID
func (cm *DockerContainerSource) Get(id string) (*Container, bool) {
	for _, c := range cm.containers {
		if c.id == id {
			return c, true
		}
	}
	return nil, false
}

// Remove containers by ID
func (cm *DockerContainerSource) delByID(id string) {
	for n, c := range cm.containers {
		if c.id == id {
			cm.del(n)
			return
		}
	}
}

// Remove one or more containers by index
func (cm *DockerContainerSource) del(idx ...int) {
	lock.Lock()
	defer lock.Unlock()
	for _, i := range idx {
		cm.containers = append(cm.containers[:i], cm.containers[i+1:]...)
	}
	log.Infof("removed %d dead containers", len(idx))
}

// Return array of all containers, sorted by field
func (cm *DockerContainerSource) All() []*Container {
	sort.Sort(cm.containers)
	return cm.containers
}
