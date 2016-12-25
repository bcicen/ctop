package main

import (
	"os"

	"github.com/fsouza/go-dockerclient"
)

func runningCIDs(client *docker.Client) (running []string) {
	filters := make(map[string][]string)
	filters["status"] = []string{"running"}
	opts := docker.ListContainersOptions{
		Filters: filters,
	}
	containers, err := client.ListContainers(opts)
	if err != nil {
		panic(err)
	}
	for _, c := range containers {
		running = append(running, c.ID[:12])
	}
	return running
}

func main() {
	var containers []string

	dockerhost := os.Getenv("DOCKER_HOST")
	if dockerhost == "" {
		dockerhost = "unix:///var/run/docker.sock"
	}

	client, err := docker.NewClient(dockerhost)
	if err != nil {
		panic(err)
	}

	// Default to all running containers
	if len(os.Args) < 2 {
		containers = runningCIDs(client)
	} else {
		containers = os.Args[1:]
	}

	g := &Grid{make(map[string]*Container)}
	for _, c := range containers {
		g.AddContainer(c)
	}

	for _, c := range g.containers {
		c.Collect(client)
	}

	Display(g)

}
