package main

import (
	"os"

	"github.com/fsouza/go-dockerclient"
)

func getContainers(client *docker.Client) []docker.APIContainers {
	filters := make(map[string][]string)
	filters["status"] = []string{"running"}
	opts := docker.ListContainersOptions{
		Filters: filters,
	}
	containers, err := client.ListContainers(opts)
	if err != nil {
		panic(err)
	}
	return containers
}

func main() {
	dockerhost := os.Getenv("DOCKER_HOST")
	if dockerhost == "" {
		dockerhost = "unix:///var/run/docker.sock"
	}

	client, err := docker.NewClient(dockerhost)
	if err != nil {
		panic(err)
	}

	g := &Grid{0, make(map[string]*Container)}
	for _, c := range getContainers(client) {
		g.AddContainer(c.ID[:12], c.Names)
	}

	for _, c := range g.containers {
		c.Collect(client)
	}

	Display(g)

}
