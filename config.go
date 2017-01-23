package main

import (
	"os"
)

var GlobalConfig = NewDefaultConfig()
var configChan = make(chan ConfigMsg)

type Config map[string]string

type ConfigMsg struct {
	key string
	val string
}

func updateConfig(k, v string) {
	configChan <- ConfigMsg{k, v}
}

// Toggle a boolean option
func (c Config) toggle(k string) {
	if c[k] == "0" {
		c[k] = "1"
	} else {
		c[k] = "0"
	}
}

func NewDefaultConfig() Config {
	docker := os.Getenv("DOCKER_HOST")
	if docker == "" {
		docker = "unix:///var/run/docker.sock"
	}
	config := Config{
		"dockerHost":   docker,
		"filterStr":    "",
		"sortField":    "id",
		"sortReverse":  "0",
		"enableHeader": "0",
	}
	go func() {
		for m := range configChan {
			config[m.key] = m.val
		}
	}()
	return config
}
