package manager

import (
	"fmt"
	api "github.com/fsouza/go-dockerclient"
	"io"
	"os"
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

// Do not allow to close reader (i.e. /dev/stdin which docker client tries to close after command execution)
type noClosableReader struct {
	wrappedReader io.Reader
}

func (w *noClosableReader) Read(p []byte) (n int, err error) {
	return w.wrappedReader.Read(p)
}

func (dc *Docker) Exec(cmd []string) error {
	execCmd, err := dc.client.CreateExec(api.CreateExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
		Container:    dc.id,
		Tty:          true,
	})

	if err != nil {
		return err
	}

	return dc.client.StartExec(execCmd.ID, api.StartExecOptions{
		InputStream:  &noClosableReader{os.Stdin},
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		RawTerminal:  true,
	})
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
