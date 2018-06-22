package manager

import (
	"fmt"
	api "github.com/fsouza/go-dockerclient"
)

type Docker struct {
	id     string
	client *api.Client
}

func NewDocker(client *api.Client, id string) *Docker {
	return &Docker{
		id:     id,
		client: client,
	}
}

func (dc *Docker) Start() error {
	c, err := dc.client.InspectContainer(dc.id)
	if err != nil {
		return fmt.Errorf("cannot inspect container: %v", err)
	}

	if err := dc.client.StartContainer(c.ID, c.HostConfig); err != nil {
		return fmt.Errorf("cannot start container: %v", err)
	}
	return nil
}

func (dc *Docker) Stop() error {
	if err := dc.client.StopContainer(dc.id, 3); err != nil {
		return fmt.Errorf("cannot stop container: %v", err)
	}
	return nil
}

func (dc *Docker) Remove() error {
	if err := dc.client.RemoveContainer(api.RemoveContainerOptions{ID: dc.id}); err != nil {
		return fmt.Errorf("cannot remove container: %v", err)
	}
	return nil
}

func (dc *Docker) Pause() error {
	if err := dc.client.PauseContainer(dc.id); err != nil {
		return fmt.Errorf("cannot pause container: %v", err)
	}
	return nil
}

func (dc *Docker) Unpause() error {
	if err := dc.client.UnpauseContainer(dc.id); err != nil {
		return fmt.Errorf("cannot unpause container: %v", err)
	}
	return nil
}

func (dc *Docker) Restart() error {
	if err := dc.client.RestartContainer(dc.id, 3); err != nil {
		return fmt.Errorf("cannot restart container: %v", err)
	}
	return nil
}
