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

func NewDefaultConfig() Config {
	docker := os.Getenv("DOCKER_HOST")
	if docker == "" {
		docker = "unix:///var/run/docker.sock"
	}
	config := Config{
		"dockerHost":   docker,
		"sortField":    "id",
		"enableHeader": "0",
	}
	go func() {
		for m := range configChan {
			config[m.key] = m.val
		}
	}()
	return config
}
