package main

import (
	"fmt"
	"os"

	"github.com/fsouza/go-dockerclient"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no container provided")
		os.Exit(1)
	}

	client, err := docker.NewClient("tcp://127.0.0.1:4243")
	if err != nil {
		panic(err)
	}

	g := &Grid{make(map[string]*Container)}
	for _, c := range os.Args[1:] {
		g.AddContainer(c)
	}

	for _, c := range g.containers {
		c.Collect(client)
	}

	Display(g)

}
