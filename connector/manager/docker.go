package manager

import (
	"fmt"
	"github.com/docker/docker/client"
	"context"
	"time"
	"github.com/docker/docker/api/types"
)

type Docker struct {
	id     string
	client *client.Client
}

func NewDocker(client *client.Client, id string) *Docker {
	return &Docker{
		id:     id,
		client: client,
	}
}

func (dc *Docker) Start() error {
	c, err := dc.client.ContainerInspect(context.Background(), dc.id)
	if err != nil {
		return fmt.Errorf("cannot inspect container: %v", err)
	}

	if err := dc.client.ContainerStart(
		context.Background(),
		c.ID,
		types.ContainerStartOptions{CheckpointID: c.HostConfig.ContainerIDFile},
	); err != nil {
		return fmt.Errorf("cannot start container: %v", err)
	}
	return nil
}

func (dc *Docker) Stop() error {
	duration := time.Duration(3)
	if err := dc.client.ContainerStop(context.Background(), dc.id, &duration); err != nil {
		return fmt.Errorf("cannot stop container: %v", err)
	}
	return nil
}

func (dc *Docker) Remove() error {
	if err := dc.client.ContainerRemove(context.Background(), dc.id, types.ContainerRemoveOptions{}); err != nil {
		return fmt.Errorf("cannot remove container: %v", err)
	}
	return nil
}
