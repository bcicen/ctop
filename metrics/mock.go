// +build !release

package metrics

import (
	"math/rand"
	"time"
)

// Mock collector
type Mock struct {
	Metrics
	stream     chan Metrics
	done       bool
	running    bool
	aggression int64
}

func NewMock(a int64) *Mock {
	c := &Mock{
		Metrics:    Metrics{},
		aggression: a,
	}
	c.MemLimit = 2147483648
	return c
}

func (c *Mock) Running() bool {
	return c.running
}

func (c *Mock) Start() {
	c.done = false
	c.stream = make(chan Metrics)
	go c.run()
}

func (c *Mock) Stop() {
	c.done = true
}

func (c *Mock) Stream() chan Metrics {
	return c.stream
}

func (c *Mock) run() {
	c.running = true
	rand.Seed(int64(time.Now().Nanosecond()))
	defer close(c.stream)

	for {
		c.CPUUtil += rand.Intn(2) * int(c.aggression)
		if c.CPUUtil >= 100 {
			c.CPUUtil = 0
		}

		c.NetTx += rand.Int63n(60) * c.aggression
		c.NetRx += rand.Int63n(60) * c.aggression
		c.MemUsage += rand.Int63n(c.MemLimit/16) * c.aggression
		if c.MemUsage > c.MemLimit {
			c.MemUsage = 0
		}
		c.MemPercent = round((float64(c.MemUsage) / float64(c.MemLimit)) * 100)
		c.stream <- c.Metrics
		if c.done {
			break
		}
		time.Sleep(1 * time.Second)
	}

	c.running = false
}
