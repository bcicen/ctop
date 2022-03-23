//go:build linux
// +build linux

package connector

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/bcicen/ctop/connector/collector"
	"github.com/bcicen/ctop/connector/manager"
	"github.com/bcicen/ctop/container"
	"github.com/opencontainers/runc/libcontainer"
)

func init() { enabled["runc"] = NewRunc }

type RuncOpts struct {
	root           string // runc root path
	systemdCgroups bool   // use systemd cgroups
}

func NewRuncOpts() (RuncOpts, error) {
	var opts RuncOpts
	// read runc root path
	root := os.Getenv("RUNC_ROOT")
	if root == "" {
		root = "/run/runc"
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return opts, err
	}
	opts.root = abs

	// ensure runc root path is readable
	_, err = ioutil.ReadDir(opts.root)
	if err != nil {
		return opts, err
	}

	if os.Getenv("RUNC_SYSTEMD_CGROUP") == "1" {
		opts.systemdCgroups = true
	}
	return opts, nil
}

type Runc struct {
	opts          RuncOpts
	factory       libcontainer.Factory
	containers    map[string]*container.Container
	libContainers map[string]libcontainer.Container
	closed        chan struct{}
	needsRefresh  chan string // container IDs requiring refresh
	lock          sync.RWMutex
}

func NewRunc() (Connector, error) {
	opts, err := NewRuncOpts()
	if err != nil {
		return nil, err
	}

	factory, err := libcontainer.New(opts.root)
	if err != nil {
		return nil, err
	}

	cm := &Runc{
		opts:          opts,
		factory:       factory,
		containers:    make(map[string]*container.Container),
		libContainers: make(map[string]libcontainer.Container),
		closed:        make(chan struct{}),
		lock:          sync.RWMutex{},
	}

	go func() {
		for {
			select {
			case <-cm.closed:
				return
			case <-time.After(5 * time.Second):
				cm.refreshAll()
			}
		}
	}()
	go cm.Loop()

	return cm, nil
}

func (cm *Runc) GetLibc(id string) libcontainer.Container {
	// return previously loaded container
	libc, ok := cm.libContainers[id]
	if ok {
		return libc
	}
	// load container
	libc, err := cm.factory.Load(id)
	if err != nil {
		// remove container if no longer exists
		if errors.Is(err, libcontainer.ErrNotExist) {
			cm.delByID(id)
		} else {
			log.Warningf("failed to read container: %s\n", err)
		}
		return nil
	}
	return libc
}

// update a ctop container from libcontainer
func (cm *Runc) refresh(id string) {
	libc := cm.GetLibc(id)
	if libc == nil {
		return
	}
	c := cm.MustGet(id)

	// remove container if entered destroyed state on last refresh
	// this gives adequate time for the collector to be shut down
	if c.GetMeta("state") == "destroyed" {
		cm.delByID(id)
		return
	}

	status, err := libc.Status()
	if err != nil {
		log.Warningf("failed to read status for container: %s\n", err)
	} else {
		c.SetState(status.String())
	}

	state, err := libc.State()
	if err != nil {
		log.Warningf("failed to read state for container: %s\n", err)
	} else {
		c.SetMeta("created", state.BaseState.Created.Format("Mon Jan 2 15:04:05 2006"))
	}

	conf := libc.Config()
	c.SetMeta("rootfs", conf.Rootfs)
}

// Read runc root, creating any new containers
func (cm *Runc) refreshAll() {
	list, err := ioutil.ReadDir(cm.opts.root)
	if err != nil {
		log.Errorf("%s (%T)", err.Error(), err)
		close(cm.closed)
		return
	}

	for _, i := range list {
		if i.IsDir() {
			name := i.Name()
			// attempt to load
			libc := cm.GetLibc(name)
			if libc == nil {
				continue
			}
			_ = cm.MustGet(i.Name()) // ensure container exists
		}
	}

	// queue all existing containers for refresh
	for id := range cm.containers {
		cm.needsRefresh <- id
	}
	log.Debugf("queued %d containers for refresh", len(cm.containers))
}

func (cm *Runc) Loop() {
	for id := range cm.needsRefresh {
		cm.refresh(id)
	}
}

// MustGet gets a single ctop container in the map matching libc container, creating one anew if not existing
func (cm *Runc) MustGet(id string) *container.Container {
	c, ok := cm.Get(id)
	if !ok {
		libc := cm.GetLibc(id)

		// create collector
		collector := collector.NewRunc(libc)

		// create container
		manager := manager.NewRunc()
		c = container.New(id, collector, manager)

		name := libc.ID()
		// set initial metadata
		if len(name) > 12 {
			name = name[0:12]
		}
		c.SetMeta("name", name)

		// add to map
		cm.lock.Lock()
		cm.containers[id] = c
		cm.libContainers[id] = libc
		cm.lock.Unlock()
		log.Debugf("saw new container: %s", id)
	}

	return c
}

// Remove containers by ID
func (cm *Runc) delByID(id string) {
	cm.lock.Lock()
	delete(cm.containers, id)
	delete(cm.libContainers, id)
	cm.lock.Unlock()
	log.Infof("removed dead container: %s", id)
}

// Runc implements Connector
func (cm *Runc) Wait() struct{} { return <-cm.closed }

// Runc implements Connector
func (cm *Runc) Get(id string) (*container.Container, bool) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	c, ok := cm.containers[id]
	return c, ok
}

// Runc implements Connector
func (cm *Runc) All() (containers container.Containers) {
	cm.lock.Lock()
	for _, c := range cm.containers {
		containers = append(containers, c)
	}
	containers.Sort()
	containers.Filter()
	cm.lock.Unlock()
	return containers
}
