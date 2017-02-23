package metrics

import (
	"math/rand"
	"time"
)

// Mock collector
type Mock struct {
	Metrics
	stream chan Metrics
	done   chan bool
}

func NewMock() *Mock {
	c := &Mock{
		Metrics: Metrics{},
		stream:  make(chan Metrics),
		done:    make(chan bool),
	}
	c.MemLimit = 2147483648
	go c.run()
	return c
}

func (c *Mock) run() {
	rand.Seed(int64(time.Now().Nanosecond()))
	for {
		c.CPUUtil += rand.Intn(10)
		if c.CPUUtil > 100 {
			c.CPUUtil = 0
		}
		c.CPUUtil += rand.Intn(2)
		c.NetTx += rand.Int63n(30)
		c.NetRx += rand.Int63n(30)
		c.MemUsage += rand.Int63n(c.MemLimit / 16)
		if c.MemUsage > c.MemLimit {
			c.MemUsage = 0
		}
		c.MemPercent = round((float64(c.MemUsage) / float64(c.MemLimit)) * 100)
		c.stream <- c.Metrics
		time.Sleep(1 * time.Second)
	}
}

func (c *Mock) Stream() chan Metrics {
	return c.stream
}

func (c *Mock) Stop() {
	c.done <- true
}
