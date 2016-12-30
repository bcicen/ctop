package main

import (
	"strings"

	"github.com/fsouza/go-dockerclient"
)

func NewContainerMap() *ContainerMap {
	return &ContainerMap{
		containers: make(map[string]*Container),
		sortField:  "cpu",
	}
}

type ContainerMap struct {
	containers map[string]*Container
	sortField  string
}

// Return number of containers/rows
func (cm *ContainerMap) Len() uint {
	return uint(len(cm.containers))
}

func (cm *ContainerMap) Add(c docker.APIContainers) {
	id := c.ID[:12]
	name := strings.Replace(c.Names[0], "/", "", 1) // use primary container name
	cm.containers[id] = &Container{
		id:      id,
		done:    make(chan bool),
		stats:   make(chan *docker.Stats),
		widgets: NewWidgets(cid, name),
		reader:  &StatReader{},
	}
}

// Get a single container, by ID
func (cm *ContainerMap) Get(id string) *Container {
	return cm.containers[id]
}

// Return array of all containers
func (cm *ContainerMap) All() []*Container {
	var containers []*Container
	for _, c := range cm.containers {
		containers = append(containers, c)
	}
	return containers
}
