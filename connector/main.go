package connector

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/bcicen/ctop/container"
	"github.com/bcicen/ctop/logging"
)

var (
	log     = logging.Init()
	enabled = make(map[string]ConnectorFn)
)

type ConnectorFn func() (Connector, error)

type Connector interface {
	// All returns a pre-sorted container.Containers of all discovered containers
	All() container.Containers
	// Get returns a single container.Container by ID
	Get(string) (*container.Container, bool)
	// Wait blocks until the underlying connection is lost
	Wait() struct{}
}

// ConnectorSuper provides initial connection and retry on failure for
// an undlerying Connector type
type ConnectorSuper struct {
	conn   Connector
	connFn ConnectorFn
	err    error
	lock   sync.RWMutex
}

func NewConnectorSuper(connFn ConnectorFn) *ConnectorSuper {
	cs := &ConnectorSuper{
		connFn: connFn,
		err:    fmt.Errorf("connecting..."),
	}
	go cs.loop()
	return cs
}

// Get returns the underlying Connector, or nil and an error
// if the Connector is not yet initialized or is disconnected.
func (cs *ConnectorSuper) Get() (Connector, error) {
	cs.lock.RLock()
	defer cs.lock.RUnlock()
	if cs.err != nil {
		return nil, cs.err
	}
	return cs.conn, nil
}

func (cs *ConnectorSuper) setError(err error) {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	cs.err = err
}

func (cs *ConnectorSuper) loop() {
	const interval = 3
	for {
		log.Infof("initializing connector")

		conn, err := cs.connFn()
		if err != nil {
			cs.setError(err)
			log.Errorf("failed to initialize connector: %s (%T)", err, err)
			log.Errorf("retrying in %ds", interval)
			time.Sleep(interval * time.Second)
		} else {
			cs.conn = conn
			cs.setError(nil)
			log.Infof("successfully initialized connector")

			// wait until connection closed
			cs.conn.Wait()
			cs.setError(fmt.Errorf("attempting to reconnect..."))
			log.Infof("connector closed")
		}
	}
}

// Enabled returns names for all enabled connectors on the current platform
func Enabled() (a []string) {
	for k, _ := range enabled {
		a = append(a, k)
	}
	sort.Strings(a)
	return a
}

// ByName returns a ConnectorSuper for a given name, or error if the connector
// does not exists on the current platform
func ByName(s string) (*ConnectorSuper, error) {
	if cfn, ok := enabled[s]; ok {
		return NewConnectorSuper(cfn), nil
	}
	return nil, fmt.Errorf("invalid connector type \"%s\"", s)
}
