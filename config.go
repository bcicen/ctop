package main

import (
	"os"
)

type Config struct {
	dockerHost string
	sortField  string
}

var DefaultConfig = NewDefaultConfig()

func NewDefaultConfig() Config {
	docker := os.Getenv("DOCKER_HOST")
	if docker == "" {
		docker = "unix:///var/run/docker.sock"
	}
	return Config{
		dockerHost: docker,
		sortField:  "id",
	}
}
