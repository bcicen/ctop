package logging

import (
	"time"

	"github.com/op/go-logging"
)

const (
	size = 1024
)

var (
	Log    *CTopLogger
	exited bool
	level  = logging.DEBUG
	format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
)

type CTopLogger struct {
	*logging.Logger
	backend *logging.MemoryBackend
}

func Init() *CTopLogger {
	if Log == nil {
		logging.SetFormatter(format) // setup default formatter

		Log = &CTopLogger{
			logging.MustGetLogger("ctop"),
			logging.NewMemoryBackend(size),
		}

		backendLvl := logging.AddModuleLevel(Log.backend)
		backendLvl.SetLevel(level, "")

		logging.SetBackend(backendLvl)
		Log.Notice("logger initialized")
	}
	return Log
}

func (log *CTopLogger) tail() chan string {
	stream := make(chan string)

	node := log.backend.Head()
	go func() {
		for {
			stream <- node.Record.Formatted(0)
			for {
				nnode := node.Next()
				if nnode != nil {
					node = nnode
					break
				}
				if exited {
					close(stream)
					return
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()

	return stream
}

func (log *CTopLogger) Exit() {
	exited = true
	StopServer()
}
