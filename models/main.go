package models

import "time"

type Log struct {
	Timestamp time.Time
	Message   string
}

type Meta map[string]string

func NewMeta() Meta { return make(Meta) }

func (m Meta) Get(k string) string {
	if s, ok := m[k]; ok {
		return s
	}
	return ""
}

type Metrics struct {
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

func NewMetrics() Metrics {
	return Metrics{
		CPUUtil:      -1,
		NetTx:        -1,
		NetRx:        -1,
		MemUsage:     -1,
		MemPercent:   -1,
		IOBytesRead:  -1,
		IOBytesWrite: -1,
		Pids:         -1,
	}
}
