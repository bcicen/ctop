package models

import "time"

type Log struct {
	Timestamp time.Time
	Message   string
}

type Meta map[string]string

// NewMeta returns an initialized Meta map.
// An optional series of key, values may be provided to populate the map prior to returning
func NewMeta(kvs ...string) Meta {
	m := make(Meta)

	var i int
	for i < len(kvs)-1 {
		m[kvs[i]] = kvs[i+1]
		i += 2
	}

	return m
}

func (m Meta) Get(k string) string {
	if s, ok := m[k]; ok {
		return s
	}
	return ""
}

type Metrics struct {
	NCpus        uint8
	CPUUtil      int
	NetTx        int64
	NetRx        int64
	MemLimit     int64
	MemPercent   int
	MemUsage     int64
	IOBytesRead  int64
	IOBytesWrite int64
	Pids         int
}

func (m *Metrics) Add(otherMetrics Metrics) {
	m.CPUUtil += otherMetrics.CPUUtil
	m.NetTx += otherMetrics.NetTx
	m.NetRx += otherMetrics.NetRx
	m.MemLimit += otherMetrics.MemLimit
	m.MemPercent += otherMetrics.MemPercent
	m.MemUsage += otherMetrics.MemUsage
	m.IOBytesRead += otherMetrics.IOBytesRead
	m.IOBytesWrite += otherMetrics.IOBytesWrite
	m.Pids += otherMetrics.Pids
}

func (m *Metrics) Subtract(otherMetrics Metrics) {
	m.CPUUtil -= otherMetrics.CPUUtil
	m.NetTx -= otherMetrics.NetTx
	m.NetRx -= otherMetrics.NetRx
	m.MemLimit -= otherMetrics.MemLimit
	m.MemPercent -= otherMetrics.MemPercent
	m.MemUsage -= otherMetrics.MemUsage
	m.IOBytesRead -= otherMetrics.IOBytesRead
	m.IOBytesWrite -= otherMetrics.IOBytesWrite
	m.Pids -= otherMetrics.Pids
}

func (m *Metrics) SumNet() int64 { return m.NetRx + m.NetTx }

func (m *Metrics) SumIO() int64 { return m.IOBytesRead + m.IOBytesWrite }
