package connector

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bcicen/ctop/connector/collector"
	"github.com/bcicen/ctop/connector/manager"
	"github.com/bcicen/ctop/container"
	api "github.com/fsouza/go-dockerclient"
)

func init() { enabled["docker"] = NewDocker }

type Docker struct {
	client       *api.Client
	containers   map[string]*container.Container
	needsRefresh chan string // container IDs requiring refresh
	lock         sync.RWMutex
}

func NewDocker() Connector {
	// init docker client
	client, err := api.NewClientFromEnv()
	if err != nil {
		panic(err)
	}
	cm := &Docker{
		client:       client,
		containers:   make(map[string]*container.Container),
		needsRefresh: make(chan string, 60),
		lock:         sync.RWMutex{},
	}
	go cm.Loop()
	cm.refreshAll()
	go cm.watchEvents()
	return cm
}

// Docker events watcher
func (cm *Docker) watchEvents() {
	log.Info("docker event listener starting")
	events := make(chan *api.APIEvents)
	cm.client.AddEventListener(events)

	for e := range events {
		if e.Type != "container" {
			continue
		}

		actionName := strings.Split(e.Action, ":")[0]

		switch actionName {
		case "start", "die", "pause", "unpause", "health_status":
			log.Debugf("handling docker event: action=%s id=%s", e.Action, e.ID)
			cm.needsRefresh <- e.ID
		case "destroy":
			log.Debugf("handling docker event: action=%s id=%s", e.Action, e.ID)
			cm.delByID(e.ID)
		}
	}
}

func portsFormat(ports map[api.Port][]api.PortBinding) string {
	var exposed []string
	var published []string

	for k, v := range ports {
		if len(v) == 0 {
			exposed = append(exposed, string(k))
			continue
		}
		for _, binding := range v {
			s := fmt.Sprintf("%s:%s -> %s", binding.HostIP, binding.HostPort, k)
			published = append(published, s)
		}
	}

	return strings.Join(append(exposed, published...), "\n")
}

func (cm *Docker) refresh(c *container.Container) {
	insp := cm.inspect(c.Id)
	// remove container if no longer exists
	if insp == nil {
		cm.delByID(c.Id)
		return
	}
	c.SetMeta("name", shortName(insp.Name))
	c.SetMeta("image", insp.Config.Image)
	c.SetMeta("ports", portsFormat(insp.NetworkSettings.Ports))
	c.SetMeta("created", insp.Created.Format("Mon Jan 2 15:04:05 2006"))
	c.SetMeta("health", insp.State.Health.Status)
	c.SetState(insp.State.Status)
}

func (cm *Docker) inspect(id string) *api.Container {
	c, err := cm.client.InspectContainer(id)
	if err != nil {
		if _, ok := err.(*api.NoSuchContainer); ok == false {
			log.Errorf(err.Error())
		}
	}
	return c
}

// Mark all container IDs for refresh
func (cm *Docker) refreshAll() {
	opts := api.ListContainersOptions{All: true}
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

func (cm *Docker) Loop() {
	for id := range cm.needsRefresh {
		c := cm.MustGet(id)
		cm.refresh(c)
	}
}

// Get a single container, creating one anew if not existing
func (cm *Docker) MustGet(id string) *container.Container {
	c, ok := cm.Get(id)
	// append container struct for new containers
	if !ok {
		// create collector
		collector := collector.NewDocker(cm.client, id)
		// create manager
		manager := manager.NewDocker(cm.client, id)
		// create container
		c = container.New(id, collector, manager)
		cm.lock.Lock()
		cm.containers[id] = c
		cm.lock.Unlock()
	}
	return c
}

// Get a single container, by ID
func (cm *Docker) Get(id string) (*container.Container, bool) {
	cm.lock.Lock()
	c, ok := cm.containers[id]
	cm.lock.Unlock()
	return c, ok
}

// Remove containers by ID
func (cm *Docker) delByID(id string) {
	cm.lock.Lock()
	delete(cm.containers, id)
	cm.lock.Unlock()
	log.Infof("removed dead container: %s", id)
}

// Return array of all containers, sorted by field
func (cm *Docker) All() (containers container.Containers) {
	cm.lock.Lock()
	for _, c := range cm.containers {
		containers = append(containers, c)
	}

	containers.Sort()
	containers.Filter()
	cm.lock.Unlock()
	return containers
}

// use primary container name
func shortName(name string) string {
	return strings.Replace(name, "/", "", 1)
}
