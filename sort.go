package main

import (
	"sort"
	"strings"

	"github.com/fsouza/go-dockerclient"
)

func NewContainerMap() *ContainerMap {
	return &ContainerMap{
		containers: make(map[string]*Container),
		sortField:  "id",
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
	cm.containers[id] = NewContainer(id, name)
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

// Return array of containers, sorted by field
func (cm *ContainerMap) Sorted() []*Container {
	containers := cm.All()

	switch cm.sortField {
	case "id":
		sort.Sort(ByID(containers))
	case "name":
		sort.Sort(ByName(containers))
	case "cpu":
		sort.Sort(ByCPU(containers))
	case "mem":
		sort.Sort(ByMem(containers))
	default:
		sort.Sort(ByID(containers))
	}

	return containers
}

type ByID []*Container

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].id < a[j].id }

type ByName []*Container

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].id < a[j].id }

type ByCPU []*Container

func (a ByCPU) Len() int           { return len(a) }
func (a ByCPU) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCPU) Less(i, j int) bool { return a[i].reader.CPUUtil < a[j].reader.CPUUtil }

type ByMem []*Container

func (a ByMem) Len() int           { return len(a) }
func (a ByMem) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByMem) Less(i, j int) bool { return a[i].reader.MemUsage < a[j].reader.MemUsage }
