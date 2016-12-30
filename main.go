package main

import (
	"os"
	"strings"

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

// Return primary container name
func parseName(names []string) string {
	return strings.Replace(names[0], "/", "", -1)
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

	g := &Grid{0, "id", make(map[string]*Container)}
	for _, c := range getContainers(client) {
		g.AddContainer(c.ID[:12], parseName(c.Names))
	}

	for _, c := range g.containers {
		c.Collect(client)
	}

	Display(g)

}
