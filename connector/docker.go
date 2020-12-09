package connector

import (
	"github.com/op/go-logging"
	"strings"
	"sync"
	"time"

	"github.com/bcicen/ctop/connector/collector"
	"github.com/bcicen/ctop/connector/manager"
	"github.com/bcicen/ctop/container"
	api "github.com/fsouza/go-dockerclient"
)

func init() { enabled["docker"] = NewDocker }

var actionToStatus = map[string]string{
	"start":   "running",
	"die":     "exited",
	"stop":    "exited",
	"pause":   "paused",
	"unpause": "running",
}

type StatusUpdate struct {
	Cid    string
	Field  string // "status" or "health"
	Status string
}

type Docker struct {
	client       *api.Client
	containers   map[string]*container.Container
	needsRefresh chan string // container IDs requiring refresh
	statuses     chan StatusUpdate
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
		needsRefresh: make(chan string, 60),
		statuses:     make(chan StatusUpdate, 60),
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
	go cm.LoopStatuses()
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

		actionName := e.Action
		// fast skip all exec_* events: exec_create, exec_start, exec_die
		if strings.HasPrefix(actionName, "exec_") {
			continue
		}
		// Action may have additional param i.e. "health_status: healthy"
		// We need to strip to have only action name
		sepIdx := strings.Index(actionName, ": ")
		if sepIdx != -1 {
			actionName = actionName[:sepIdx]
		}

		switch actionName {
		// most frequent event is a health checks
		case "health_status":
			healthStatus := e.Action[sepIdx+2:]
			if log.IsEnabledFor(logging.DEBUG) {
				log.Debugf("handling docker event: action=health_status id=%s %s", e.ID, healthStatus)
			}
			cm.statuses <- StatusUpdate{e.ID, "health", healthStatus}
		case "create":
			if log.IsEnabledFor(logging.DEBUG) {
				log.Debugf("handling docker event: action=create id=%s", e.ID)
			}
			c := cm.MustGet(e.ID)
			c.SetMeta("name", manager.ShortName(e.Actor.Attributes["name"]))
			c.SetMeta("image", e.Actor.Attributes["image"])
			c.SetState("created")
			cm.needsRefresh <- e.ID
		case "destroy":
			if log.IsEnabledFor(logging.DEBUG) {
				log.Debugf("handling docker event: action=destroy id=%s", e.ID)
			}
			cm.delByID(e.ID)
		default:
			// check if this action changes status e.g. start -> running
			status := actionToStatus[actionName]
			if status != "" {
				if log.IsEnabledFor(logging.DEBUG) {
					log.Debugf("handling docker event: action=%s id=%s %s", actionName, e.ID, status)
				}
				cm.statuses <- StatusUpdate{e.ID, "status", status}
			}
		}
	}
	log.Info("docker event listener exited")
	close(cm.closed)
}

func (cm *Docker) refresh(c *container.Container) {
	cm.updateContainers(false, []string{c.Id})
}

// Mark all container IDs for refresh
func (cm *Docker) refreshAll() {
	cm.updateContainers(true, nil)
}

func (cm *Docker) updateContainers(all bool, cidsToRefresh []string) {
	opts := api.ListContainersOptions{All: true}
	if !all {
		opts.Filters = map[string][]string{
			"id": cidsToRefresh,
		}
	}
	allContainers, err := cm.client.ListContainers(opts)
	if err != nil {
		log.Errorf("%s (%T)", err.Error(), err)
		return
	}
	if all {
		cm.cleanupDestroyedContainers(allContainers)
	}

	for _, i := range allContainers {
		c := cm.MustGet(i.ID)
		c.SetMeta("name", manager.ShortName(i.Names[0]))
		c.SetMeta("image", i.Image)
		c.SetMeta("IPs", manager.IpsFormat(i.Networks.Networks))
		c.SetMeta("ports", manager.PortsFormatArr(i.Ports))
		c.SetMeta("created", time.Unix(i.Created, 0).Format("Mon Jan 2 15:04:05 2006"))
		parseStatusHealth(c, i.Status)
		c.SetState(i.State)
	}
}

func (cm *Docker) cleanupDestroyedContainers(allContainers []api.APIContainers) {
	var nonExistingContainers []string
	for _, oldContainer := range cm.containers {
		if !cm.hasContainer(oldContainer.Id, allContainers) {
			nonExistingContainers = append(nonExistingContainers, oldContainer.Id)
		}
	}
	// remove containers that no longer exists
	for _, cid := range nonExistingContainers {
		cm.delByID(cid)
	}
}

func (cm *Docker) hasContainer(oldContainerId string, newContainers []api.APIContainers) bool {
	for _, newContainer := range newContainers {
		if newContainer.ID == oldContainerId {
			return true
		}
	}
	return false
}

func parseStatusHealth(c *container.Container, status string) {
	// Status may look like:
	//  Up About a minute (healthy)
	//  Up 7 minutes (unhealthy)
	var health string
	if strings.Contains(status, "(healthy)") {
		health = "healthy"
	} else if strings.Contains(status, "(unhealthy)") {
		health = "unhealthy"
	} else {
		return
	}
	c.SetMeta("health", health)
}

func (cm *Docker) Loop() {
	ticker := time.NewTicker(5 * time.Minute)
	for {
		select {
		case <-ticker.C:
			cm.refreshAll()
		case id := <-cm.needsRefresh:
			c := cm.MustGet(id)
			cm.refresh(c)
		case <-cm.closed:
			ticker.Stop()
			return
		}
	}
}

func (cm *Docker) LoopStatuses() {
	for {
		select {
		case statusUpdate := <-cm.statuses:
			c, _ := cm.Get(statusUpdate.Cid)
			if c != nil {
				if statusUpdate.Field == "health" {
					c.SetMeta("health", statusUpdate.Status)
				} else {
					c.SetState(statusUpdate.Status)
				}
			}
		case <-cm.closed:
			return
		}
	}
}

// MustGet gets a single container, creating one anew if not existing
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
	delete(cm.containers, id)
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
