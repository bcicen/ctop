package manager

import (
	"fmt"
	"github.com/bcicen/ctop/logging"
	"github.com/bcicen/ctop/models"
	api "github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
	"io"
	"os"
)

var (
	log = logging.Init()
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
	io.Reader
}

func (w *noClosableReader) Read(p []byte) (n int, err error) {
	return w.Reader.Read(p)
}

const (
	STDIN  = 0
	STDOUT = 1
	STDERR = 2
)

var wrongFrameFormat = errors.New("Wrong frame format")

// A frame has a Header and a Payload
// Header: [8]byte{STREAM_TYPE, 0, 0, 0, SIZE1, SIZE2, SIZE3, SIZE4}
// STREAM_TYPE can be:
//    0: stdin (is written on stdout)
//    1: stdout
//    2: stderr
// SIZE1, SIZE2, SIZE3, SIZE4 are the four bytes of the uint32 size encoded as big endian.
// But we don't use size, because we don't need to find the end of frame.
type frameWriter struct {
	stdout io.Writer
	stderr io.Writer
	stdin  io.Writer
}

func (w *frameWriter) Write(p []byte) (n int, err error) {
	// drop initial empty frames
	if len(p) == 0 {
		return 0, nil
	}

	if len(p) > 8 {
		var targetWriter io.Writer
		switch p[0] {
		case STDIN:
			targetWriter = w.stdin
			break
		case STDOUT:
			targetWriter = w.stdout
			break
		case STDERR:
			targetWriter = w.stderr
			break
		default:
			return 0, wrongFrameFormat
		}

		n, err := targetWriter.Write(p[8:])
		return n + 8, err
	}

	return 0, wrongFrameFormat
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
		OutputStream: &frameWriter{os.Stdout, os.Stderr, os.Stdin},
		ErrorStream:  os.Stderr,
		RawTerminal:  true,
	})
}

func (dc *Docker) inspect(id string) (insp *api.Container, found bool, err error) {
	c, err := dc.client.InspectContainer(id)
	if err != nil {
		if _, notFound := err.(*api.NoSuchContainer); notFound {
			return c, false, nil
		}
		// other error e.g. connection failed
		log.Errorf("%s (%T)", err.Error(), err)
		return c, false, err
	}
	return c, true, nil
}

func (dc *Docker) Inspect() (models.Meta, error) {
	insp, found, err := dc.inspect(dc.id)
	if !found {
		return nil, err
	}
	newMeta := models.Meta{}
	newMeta["name"] = ShortName(insp.Name)
	newMeta["image"] = insp.Config.Image
	newMeta["IPs"] = IpsFormat(insp.NetworkSettings.Networks)
	newMeta["ports"] = PortsFormat(insp.NetworkSettings.Ports)
	newMeta["created"] = insp.Created.Format("Mon Jan 2 15:04:05 2006")
	newMeta["health"] = insp.State.Health.Status
	for _, env := range insp.Config.Env {
		newMeta["[ENV-VAR]"] = env
	}
	newMeta["state"] = insp.State.Status

	return newMeta, nil
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
