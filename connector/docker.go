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
	noneProject  *container.Project
	projects     map[string]*container.Project
	needsRefresh chan string // container IDs requiring refresh
	closed       chan struct{}
	lock         sync.RWMutex
}

func NewDocker() (Connector, error) {
	// init docker client
	client, err := api.NewClientFromEnv()
	if err != nil {
		return nil, err
	}
	cm := &Docker{
		client:       client,
		containers:   make(map[string]*container.Container),
		noneProject:  container.NewProject(""),
		projects:     make(map[string]*container.Project),
		needsRefresh: make(chan string, 60),
		closed:       make(chan struct{}),
		lock:         sync.RWMutex{},
	}

	// query info as pre-flight healthcheck
	info, err := client.Info()
	if err != nil {
		return nil, err
	}

	log.Debugf("docker-connector ID: %s", info.ID)
	log.Debugf("docker-connector Driver: %s", info.Driver)
	log.Debugf("docker-connector Images: %d", info.Images)
	log.Debugf("docker-connector Name: %s", info.Name)
	log.Debugf("docker-connector ServerVersion: %s", info.ServerVersion)

	go cm.Loop()
	cm.refreshAll()
	go cm.watchEvents()
	return cm, nil
}

// Docker implements Connector
func (cm *Docker) Wait() struct{} { return <-cm.closed }

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
	log.Info("docker event listener exited")
	close(cm.closed)
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

func ipsFormat(networks map[string]api.ContainerNetwork) string {
	var ips []string

	for k, v := range networks {
		s := fmt.Sprintf("%s:%s", k, v.IPAddress)
		ips = append(ips, s)
	}

	return strings.Join(ips, "\n")
}

func (cm *Docker) refresh(id string) {
	insp := cm.inspect(id)
	// remove container if no longer exists
	if insp == nil {
		cm.delByID(id)
		return
	}
	c := cm.MustGet(id, insp.Config.Labels)
	c.SetMeta("name", shortName(insp.Name))
	c.SetMeta("image", insp.Config.Image)
	c.SetMeta("IPs", ipsFormat(insp.NetworkSettings.Networks))
	c.SetMeta("ports", portsFormat(insp.NetworkSettings.Ports))
	c.SetMeta("created", insp.Created.Format("Mon Jan 2 15:04:05 2006"))
	c.SetMeta("health", insp.State.Health.Status)
	for _, env := range insp.Config.Env {
		c.SetMeta("[ENV-VAR]", env)
	}
	c.SetState(insp.State.Status)
}

func (cm *Docker) inspect(id string) *api.Container {
	c, err := cm.client.InspectContainer(id)
	if err != nil {
		if _, ok := err.(*api.NoSuchContainer); !ok {
			log.Errorf("%s (%T)", err.Error(), err)
		}
	}
	return c
}

// Mark all container IDs for refresh
func (cm *Docker) refreshAll() {
	opts := api.ListContainersOptions{All: true}
	allContainers, err := cm.client.ListContainers(opts)
	if err != nil {
		log.Errorf("%s (%T)", err.Error(), err)
		return
	}

	for _, i := range allContainers {
		c := cm.MustGet(i.ID, i.Labels)
		c.SetMeta("name", shortName(i.Names[0]))
		c.SetState(i.State)
		cm.needsRefresh <- c.Id
	}
}

func (cm *Docker) initContainerProject(c *container.Container, labels map[string]string) {
	projectName := labels["com.docker.compose.project"]
	if projectName == "" {
		c.Project = cm.noneProject
	} else {
		// try to find the existing project
		project := cm.projects[projectName]
		if project != nil {
			c.Project = project
		} else {
			// create and remember the new project
			c.Project = container.NewProject(projectName)
			c.Project.WorkDir = labels["com.docker.compose.project.working_dir"]
			c.Project.Config = labels["com.docker.compose.project.config_files"]
			cm.projects[projectName] = c.Project
			// set compose service for the container
			composeService := labels["com.docker.compose.service"]
			if composeService != "" {
				c.SetMeta("service", composeService)
			}
		}
	}
	c.Project.Count++
}

func (cm *Docker) Loop() {
	for {
		select {
		case id := <-cm.needsRefresh:
			cm.refresh(id)
		case <-cm.closed:
			return
		}
	}
}

// MustGet gets a single container, creating one anew if not existing
func (cm *Docker) MustGet(id string, labels map[string]string) *container.Container {
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
		cm.initContainerProject(c, labels)
		cm.lock.Unlock()
	}
	return c
}

// Docker implements Connector
func (cm *Docker) Get(id string) (*container.Container, bool) {
	cm.lock.Lock()
	c, ok := cm.containers[id]
	cm.lock.Unlock()
	return c, ok
}

// Remove containers by ID
func (cm *Docker) delByID(id string) {
	cm.lock.Lock()
	c, hasContainer := cm.containers[id]
	if hasContainer {
		c.Project.Count--
		// if this was the last container in project then remove project
		if c.Project != cm.noneProject && c.Project.Count <= 0 {
			delete(cm.projects, c.Project.Name)
		}
		delete(cm.containers, id)
	}
	cm.lock.Unlock()
	log.Infof("removed dead container: %s", id)
}

// Docker implements Connector
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
