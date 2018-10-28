package manager

import (
	"k8s.io/client-go/kubernetes"
)

type Kubernetes struct {
	id     string
	client *kubernetes.Clientset
}

func NewKubernetes(client *kubernetes.Clientset, id string) *Kubernetes {
	return &Kubernetes{
		id:     id,
		client: client,
	}
}

func (dc *Kubernetes) Start() error {
	//c, err := dc.client.InspectContainer(dc.id)
	//if err != nil {
	//	return fmt.Errorf("cannot inspect container: %v", err)
	//}

	//if err := dc.client.StartContainer(c.ID, c.HostConfig); err != nil {
	//	return fmt.Errorf("cannot start container: %v", err)
	//}
	return nil
}

func (dc *Kubernetes) Stop() error {
	//if err := dc.client.StopContainer(dc.id, 3); err != nil {
	//	return fmt.Errorf("cannot stop container: %v", err)
	//}
	return nil
}

func (dc *Kubernetes) Remove() error {
	//if err := dc.client.RemoveContainer(api.RemoveContainerOptions{ID: dc.id}); err != nil {
	//	return fmt.Errorf("cannot remove container: %v", err)
	//}
	return nil
}

func (dc *Kubernetes) Pause() error {
	//if err := dc.client.PauseContainer(dc.id); err != nil {
	//	return fmt.Errorf("cannot pause container: %v", err)
	//}
	return nil
}

func (dc *Kubernetes) Unpause() error {
	//if err := dc.client.UnpauseContainer(dc.id); err != nil {
	//	return fmt.Errorf("cannot unpause container: %v", err)
	//}
	return nil
}

func (dc *Kubernetes) Restart() error {
	//if err := dc.client.RestartContainer(dc.id, 3); err != nil {
	//	return fmt.Errorf("cannot restart container: %v", err)
	//}
	return nil
}
