package main

import (
	"sort"

	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/metrics"
	"github.com/fsouza/go-dockerclient"
)

type ContainerMap struct {
	client       *docker.Client
	containers   Containers
	collectors   map[string]metrics.Collector
	needsRefresh map[string]int // container IDs requiring refresh
}

func NewContainerMap() *ContainerMap {
	// init docker client
	client, err := docker.NewClient(config.GetVal("dockerHost"))
	if err != nil {
		panic(err)
	}
	cm := &ContainerMap{
		client:       client,
		collectors:   make(map[string]metrics.Collector),
		needsRefresh: make(map[string]int),
	}
	cm.refreshAll()
	go cm.watch()
	return cm
}

// Docker events watcher
func (cm *ContainerMap) watch() {
	log.Info("docker event listener starting")
	events := make(chan *docker.APIEvents)
	cm.client.AddEventListener(events)

	for e := range events {
		cm.handleEvent(e)
	}
}

// Docker event handler
func (cm *ContainerMap) handleEvent(e *docker.APIEvents) {
	// only process container events
	if e.Type != "container" {
		return
	}
	switch e.Action {
	case "start", "die", "pause", "unpause":
		cm.needsRefresh[e.ID] = 1
	case "destroy":
		cm.DelByID(e.ID)
	}
}

func (cm *ContainerMap) refresh(id string) {
	insp := cm.inspect(id)
	// remove container if no longer exists
	if insp == nil {
		cm.DelByID(id)
		return
	}

	c, ok := cm.Get(id)
	// append container struct for new containers
	if !ok {
		c = &Container{
			id:   id,
			name: insp.Name,
		}
		c.Collapse()
		cm.containers = append(cm.containers, c)
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
}

func (cm *ContainerMap) inspect(id string) *docker.Container {
	c, err := cm.client.InspectContainer(id)
	if err != nil {
		if _, ok := err.(*docker.NoSuchContainer); ok == false {
			log.Errorf(err.Error())
		}
	}
	return c
}

func (cm *ContainerMap) refreshAll() {
	opts := docker.ListContainersOptions{All: true}
	allContainers, err := cm.client.ListContainers(opts)
	if err != nil {
		panic(err)
	}

	for _, c := range allContainers {
		cm.needsRefresh[c.ID] = 1
	}
	cm.Update()
}

func (cm *ContainerMap) Update() {
	var ids []string
	for id, _ := range cm.needsRefresh {
		cm.refresh(id)
		ids = append(ids, id)
	}
	for _, id := range ids {
		delete(cm.needsRefresh, id)
	}
}

// Kill a container by ID
//func (cm *ContainerMap) Kill(id string, sig docker.Signal) error {
//opts := docker.KillContainerOptions{
//ID:     id,
//Signal: sig,
//}
//return cm.client.KillContainer(opts)
//}

// Return number of containers/rows
func (cm *ContainerMap) Len() uint {
	return uint(len(cm.containers))
}

// Get a single container, by ID
func (cm *ContainerMap) Get(id string) (*Container, bool) {
	for _, c := range cm.containers {
		if c.id == id {
			return c, true
		}
	}
	return nil, false
}

// Remove containers by ID
func (cm *ContainerMap) DelByID(id string) {
	for n, c := range cm.containers {
		if c.id == id {
			cm.Del(n)
			return
		}
	}
}

// Remove one or more containers by index
func (cm *ContainerMap) Del(idx ...int) {
	for _, i := range idx {
		cm.containers = append(cm.containers[:i], cm.containers[i+1:]...)
	}
}

// Return array of all containers, sorted by field
func (cm *ContainerMap) All() []*Container {
	sort.Sort(cm.containers)
	return cm.containers
}
