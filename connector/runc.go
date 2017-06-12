package connector

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/bcicen/ctop/container"
	"github.com/bcicen/ctop/metrics"
	"github.com/opencontainers/runc/libcontainer"
	"github.com/opencontainers/runc/libcontainer/cgroups/systemd"
)

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
	needsRefresh  chan string // container IDs requiring refresh
	lock          sync.RWMutex
}

func NewRunc() Connector {
	opts, err := NewRuncOpts()
	runcFailOnErr(err)

	factory, err := getFactory(opts)
	runcFailOnErr(err)

	cm := &Runc{
		opts:          opts,
		factory:       factory,
		containers:    make(map[string]*container.Container),
		libContainers: make(map[string]libcontainer.Container),
		needsRefresh:  make(chan string, 60),
		lock:          sync.RWMutex{},
	}

	go func() {
		for {
			cm.refreshAll()
			time.Sleep(5 * time.Second)
		}
	}()
	go cm.Loop()

	return cm
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
		if lerr, ok := err.(libcontainer.Error); ok && lerr.Code() == libcontainer.ContainerNotExists {
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
	runcFailOnErr(err)

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
	for id, _ := range cm.containers {
		cm.needsRefresh <- id
	}
	log.Debugf("queued %d containers for refresh", len(cm.containers))
}

func (cm *Runc) Loop() {
	for id := range cm.needsRefresh {
		cm.refresh(id)
	}
}

// Get a single ctop container in the map matching libc container, creating one anew if not existing
func (cm *Runc) MustGet(id string) *container.Container {
	c, ok := cm.Get(id)
	if !ok {
		libc := cm.GetLibc(id)

		// create collector
		collector := metrics.NewRunc(libc)

		// create container
		c = container.New(id, collector)

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

// Get a single container, by ID
func (cm *Runc) Get(id string) (*container.Container, bool) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	c, ok := cm.containers[id]
	return c, ok
}

// Remove containers by ID
func (cm *Runc) delByID(id string) {
	cm.lock.Lock()
	delete(cm.containers, id)
	delete(cm.libContainers, id)
	cm.lock.Unlock()
	log.Infof("removed dead container: %s", id)
}

// Return array of all containers, sorted by field
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

func getFactory(opts RuncOpts) (libcontainer.Factory, error) {
	cgroupManager := libcontainer.Cgroupfs
	if opts.systemdCgroups {
		if systemd.UseSystemd() {
			cgroupManager = libcontainer.SystemdCgroups
		} else {
			return nil, fmt.Errorf("systemd cgroup enabled, but systemd support for managing cgroups is not available")
		}
	}
	return libcontainer.New(opts.root, cgroupManager)
}

func runcFailOnErr(err error) {
	if err != nil {
		panic(fmt.Errorf("fatal runc error: %s", err))
	}
}
