package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/bcicen/ctop/collector"
	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/widgets"
	"github.com/fsouza/go-dockerclient"
)

func NewContainerMap() *ContainerMap {
	// init docker client
	client, err := docker.NewClient(config.Get("dockerHost"))
	if err != nil {
		panic(err)
	}
	cm := &ContainerMap{
		client:     client,
		containers: make(map[string]*Container),
	}
	cm.Refresh()
	return cm
}

type ContainerMap struct {
	client     *docker.Client
	containers map[string]*Container
}

func (cm *ContainerMap) Refresh() {
	var id, name string

	opts := docker.ListContainersOptions{All: true}
	containers, err := cm.client.ListContainers(opts)
	if err != nil {
		panic(err)
	}

	// add new containers
	states := make(map[string]string)
	for _, c := range containers {
		id = c.ID[:12]
		states[id] = c.State
		if _, ok := cm.containers[id]; ok == false {
			name = strings.Replace(c.Names[0], "/", "", 1) // use primary container name
			cm.containers[id] = &Container{
				id:      id,
				name:    name,
				collect: collector.NewDocker(cm.client, id),
				widgets: widgets.NewCompact(id, name),
			}
		}
	}

	var removeIDs []string
	for id, c := range cm.containers {
		// mark stale internal containers
		if _, ok := states[id]; ok == false {
			removeIDs = append(removeIDs, id)
			continue
		}
		c.SetState(states[id])
	}

	// remove dead containers
	cm.Del(removeIDs...)
}

// Kill a container by ID
func (cm *ContainerMap) Kill(id string, sig docker.Signal) error {
	opts := docker.KillContainerOptions{
		ID:     id,
		Signal: sig,
	}
	return cm.client.KillContainer(opts)
}

// Return number of containers/rows
func (cm *ContainerMap) Len() uint {
	return uint(len(cm.containers))
}

// Get a single container, by ID
func (cm *ContainerMap) Get(id string) *Container {
	return cm.containers[id]
}

// Remove one or more containers
func (cm *ContainerMap) Del(ids ...string) {
	for _, id := range ids {
		delete(cm.containers, id)
	}
}

// Return array of all containers, sorted by field
func (cm *ContainerMap) All() []*Container {
	var containers Containers

	filter := config.Get("filterStr")
	re := regexp.MustCompile(fmt.Sprintf(".*%s", filter))

	for _, c := range cm.containers {
		if re.FindAllString(c.name, 1) != nil {
			containers = append(containers, c)
		}
	}

	sort.Sort(containers)
	return containers
}
