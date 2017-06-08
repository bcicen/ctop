package connector

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
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

type Runc struct {
	opts         RuncOpts
	factory      libcontainer.Factory
	containers   map[string]*container.Container
	needsRefresh chan string // container IDs requiring refresh
	lock         sync.RWMutex
}

func readRuncOpts() (RuncOpts, error) {
	var opts RuncOpts
	// read runc root path
	root := os.Getenv("RUNC_ROOT")
	if root == "" {
		return opts, fmt.Errorf("RUNC_ROOT not set")
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return opts, err
	}
	opts.root = abs

	if os.Getenv("RUNC_SYSTEMD_CGROUP") == "1" {
		opts.systemdCgroups = true
	}
	return opts, nil
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
	return libcontainer.New(opts.root, cgroupManager, libcontainer.CriuPath("criu"))
}

func NewRunc() *Runc {
	opts, err := readRuncOpts()
	runcFailOnErr(err)

	factory, err := getFactory(opts)
	runcFailOnErr(err)

	cm := &Runc{
		opts:         opts,
		factory:      factory,
		containers:   make(map[string]*container.Container),
		needsRefresh: make(chan string, 60),
		lock:         sync.RWMutex{},
	}

	go cm.Loop()
	go func() {
		time.Sleep(1 * time.Second)
		for _, c := range cm.containers {
			cm.needsRefresh <- c.Id
		}
	}()

	return cm
}

func (cm *Runc) inspect(id string) libcontainer.Container {
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

func (cm *Runc) refresh(c *container.Container) {
	libc := cm.inspect(c.Id)
	if libc == nil {
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

func (cm *Runc) refreshAll() {
	list, err := ioutil.ReadDir(cm.opts.root)
	runcFailOnErr(err)

	for _, i := range list {
		if i.IsDir() {
			// attempt to load
			libc, err := cm.factory.Load(i.Name())
			if err != nil {
				log.Warningf("failed to read container: %s\n", err)
				continue
			}

			c := cm.MustGet(libc.ID())
			c.SetMeta("name", i.Name())

			cm.needsRefresh <- c.Id
		}
	}
}

func (cm *Runc) Loop() {
	for id := range cm.needsRefresh {
		c := cm.MustGet(id)
		cm.refresh(c)
	}
}

// Get a single container in the map, creating one anew if not existing
func (cm *Runc) MustGet(id string) *container.Container {
	c, ok := cm.Get(id)
	// append container struct for new containers
	if !ok {
		// create collector
		collector := metrics.NewRunc(2)
		// create container
		c = container.New(id, collector)
		cm.lock.Lock()
		cm.containers[id] = c
		cm.lock.Unlock()
	}
	return c
}

// Get a single container, by ID
func (cm *Runc) Get(id string) (*container.Container, bool) {
	cm.lock.Lock()
	c, ok := cm.containers[id]
	cm.lock.Unlock()
	return c, ok
}

// Remove containers by ID
func (cm *Runc) delByID(id string) {
	cm.lock.Lock()
	delete(cm.containers, id)
	cm.lock.Unlock()
	log.Infof("removed dead container: %s", id)
}

// Return array of all containers, sorted by field
func (cm *Runc) All() (containers container.Containers) {
	cm.lock.Lock()
	for _, c := range cm.containers {
		containers = append(containers, c)
	}
	cm.lock.Unlock()
	sort.Sort(containers)
	containers.Filter()
	return containers
}

func runcFailOnErr(err error) {
	if err != nil {
		panic(fmt.Errorf("fatal runc error: %s", err))
	}
}
