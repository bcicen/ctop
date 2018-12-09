package models

import "time"

type Log struct {
	Timestamp time.Time
	Message   string
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
	Uptime       int64
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
		Uptime:       -1,
	}
}
