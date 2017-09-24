package models

import "time"

type Log struct {
	Timestamp time.Time
	Message   string
}

type Metrics struct {
	Id           string `json:"id"`
	CPUUtil      int    `json:"cpu_util"`
	NetTx        int64  `json:"net_tx"`
	NetRx        int64  `json:"net_rx"`
	MemLimit     int64  `json:"mem_limit"`
	MemPercent   int    `json:"mem_percent"`
	MemUsage     int64  `json:"mem_usage"`
	IOBytesRead  int64  `json:"io_bytes_read"`
	IOBytesWrite int64  `json:"io_bytes_write"`
	Pids         int    `json:"pids"`
}

func NewMetrics() Metrics {
	return Metrics{
		Id:           "",
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
