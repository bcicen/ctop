package main

import (
	"sort"
	"strings"

	"github.com/bcicen/ctop/metrics"
	"github.com/fsouza/go-dockerclient"
)

type ContainerSource interface {
	All() Containers
	Get(string) (*Container, bool)
}

type DockerContainerSource struct {
	client       *docker.Client
	containers   map[string]*Container
	needsRefresh chan string // container IDs requiring refresh
}

func NewDockerContainerSource() *DockerContainerSource {
	// init docker client
	client, err := docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}
	cm := &DockerContainerSource{
		client:       client,
		containers:   make(map[string]*Container),
		needsRefresh: make(chan string, 60),
	}
	go cm.Loop()
	cm.refreshAll()
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
			cm.needsRefresh <- e.ID
		case "destroy":
			log.Debugf("handling docker event: action=%s id=%s", e.Action, e.ID)
			cm.delByID(e.ID)
		}
	}
}

func (cm *DockerContainerSource) refresh(c *Container) {
	insp := cm.inspect(c.Id)
	// remove container if no longer exists
	if insp == nil {
		cm.delByID(c.Id)
		return
	}
	c.SetMeta("name", shortName(insp.Name))
	c.SetMeta("image", insp.Config.Image)
	c.SetMeta("created", insp.Created.Format("Mon Jan 2 15:04:05 2006"))
	c.SetState(insp.State.Status)
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

	for _, i := range allContainers {
		c := cm.MustGet(i.ID)
		c.SetMeta("name", shortName(i.Names[0]))
		c.SetState(i.State)
		cm.needsRefresh <- c.Id
	}
}

func (cm *DockerContainerSource) Loop() {
	for id := range cm.needsRefresh {
		c := cm.MustGet(id)
		cm.refresh(c)
	}
}

// Get a single container, creating one anew if not existing
func (cm *DockerContainerSource) MustGet(id string) *Container {
	c, ok := cm.Get(id)
	// append container struct for new containers
	if !ok {
		// create collector
		collector := metrics.NewDocker(cm.client, id)
		// create container
		c = NewContainer(id, collector)
		cm.containers[id] = c
	}
	return c
}

// Get a single container, by ID
func (cm *DockerContainerSource) Get(id string) (*Container, bool) {
	c, ok := cm.containers[id]
	return c, ok
}

// Remove containers by ID
func (cm *DockerContainerSource) delByID(id string) {
	delete(cm.containers, id)
	log.Infof("removed dead container: %s", id)
}

// Return array of all containers, sorted by field
func (cm *DockerContainerSource) All() (containers Containers) {
	for _, c := range cm.containers {
		containers = append(containers, c)
	}
	sort.Sort(containers)
	containers.Filter()
	return containers
}

// use primary container name
func shortName(name string) string {
	return strings.Replace(name, "/", "", 1)
}
