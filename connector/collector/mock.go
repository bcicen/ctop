// +build !release

package collector

import (
	"math/rand"
	"time"

	"github.com/bcicen/ctop/models"
)

const mockLog = "Cura ob pro qui tibi inveni dum qua fit donec amare illic mea, regem falli contexo pro peregrinorum heremo absconditi araneae meminerim deliciosas actionibus facere modico dura sonuerunt psalmi contra rerum, tempus mala anima volebant dura quae o modis."

// Mock collector
type Mock struct {
	models.Metrics
	stream     chan models.Metrics
	done       bool
	running    bool
	aggression int64
}

func NewMock(a int64) *Mock {
	c := &Mock{
		Metrics:    models.Metrics{},
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
	c.stream = make(chan models.Metrics)
	go c.run()
}

func (c *Mock) Stop() {
	c.done = true
}

func (c *Mock) Stream() chan models.Metrics {
	return c.stream
}

func (c *Mock) StreamLogs() (chan string, error) {
	logCh := make(chan string)
	go func() {
		for {
			logCh <- mockLog
			time.Sleep(250 * time.Millisecond)
		}
	}()
	return logCh, nil
}

func (c *Mock) run() {
	c.running = true
	rand.Seed(int64(time.Now().Nanosecond()))
	defer close(c.stream)

	// set to random static value, once
	c.Pids = rand.Intn(12)
	c.IOBytesRead = rand.Int63n(8098) * c.aggression
	c.IOBytesWrite = rand.Int63n(8098) * c.aggression

	for {
		c.CPUUtil += rand.Intn(2) * int(c.aggression)
		if c.CPUUtil >= 100 {
			c.CPUUtil = 0
		}

		c.NetTx += rand.Int63n(60) * c.aggression
		c.NetRx += rand.Int63n(60) * c.aggression
		c.MemUsage += rand.Int63n(c.MemLimit/512) * c.aggression
		if c.MemUsage > c.MemLimit {
			c.MemUsage = 0
		}
		c.MemPercent = percent(float64(c.MemUsage), float64(c.MemLimit))
		c.stream <- c.Metrics
		if c.done {
			break
		}
		time.Sleep(1 * time.Second)
	}

	c.running = false
}
