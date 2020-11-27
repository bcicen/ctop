package manager

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"io"
	"time"
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
	ctx := context.Background()
	execConfig := types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
		Tty:          true,
	}
	execCmd, err := dc.client.ContainerExecCreate(ctx, dc.id, execConfig)
	if err != nil {
		return err
	}

	execStartConfig := types.ExecStartCheck{}
	//execStartConfig := types.ExecStartCheck{
	//	InputStream:  &noClosableReader{os.Stdin},
	//	OutputStream: &frameWriter{os.Stdout, os.Stderr, os.Stdin},
	//	ErrorStream:  os.Stderr,
	//	RawTerminal:  true,
	//}
	return dc.client.ContainerExecStart(ctx, execCmd.ID, execStartConfig)
}

func (dc *Docker) Start() error {
	ctx := context.Background()
	opts := types.ContainerStartOptions{}
	if err := dc.client.ContainerStart(ctx, dc.id, opts); err != nil {
		return fmt.Errorf("cannot start container: %v", err)
	}
	return nil
}

func (dc *Docker) Stop() error {
	ctx := context.Background()
	timeout := 3 * time.Second
	if err := dc.client.ContainerStop(ctx, dc.id, &timeout); err != nil {
		return fmt.Errorf("cannot stop container: %v", err)
	}
	return nil
}

func (dc *Docker) Remove() error {
	ctx := context.Background()
	opts := types.ContainerRemoveOptions{}
	if err := dc.client.ContainerRemove(ctx, dc.id, opts); err != nil {
		return fmt.Errorf("cannot remove container: %v", err)
	}
	return nil
}

func (dc *Docker) Pause() error {
	ctx := context.Background()
	if err := dc.client.ContainerPause(ctx, dc.id); err != nil {
		return fmt.Errorf("cannot pause container: %v", err)
	}
	return nil
}

func (dc *Docker) Unpause() error {
	ctx := context.Background()
	if err := dc.client.ContainerUnpause(ctx, dc.id); err != nil {
		return fmt.Errorf("cannot unpause container: %v", err)
	}
	return nil
}

func (dc *Docker) Restart() error {
	ctx := context.Background()
	timeout := 3 * time.Second
	if err := dc.client.ContainerRestart(ctx, dc.id, &timeout); err != nil {
		return fmt.Errorf("cannot restart container: %v", err)
	}
	return nil
}
