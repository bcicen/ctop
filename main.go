package main

import (
	"fmt"
	"os"

	"github.com/fsouza/go-dockerclient"
)

type DTop struct {
	client *docker.Client
	stats  chan *docker.Stats
}

func (dt *DTop) output() {
	for s := range dt.stats {
		fmt.Println(s)
	}
}

func (dt *DTop) collect(containerID string) {
	done := make(chan bool)

	fmt.Sprintf("starting collector for container: %s\n", containerID)
	opts := docker.StatsOptions{
		ID:     containerID,
		Stats:  dt.stats,
		Stream: true,
		Done:   done,
	}
	dt.client.Stats(opts)
	fmt.Sprintf("stopping collector for container: %s\n", containerID)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no container provided")
		os.Exit(1)
	}

	client, err := docker.NewClient("tcp://127.0.0.1:4243")
	if err != nil {
		panic(err)
	}

	d := &DTop{client, make(chan *docker.Stats)}
	go d.collect(os.Args[1])
	d.output()
}
