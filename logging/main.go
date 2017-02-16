package logging

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/op/go-logging"
)

const (
	size = 1024
	path = "/tmp/ctop.sock"
)

var (
	Log    *CTopLogger
	wg     sync.WaitGroup
	exited bool
	level  = logging.INFO
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

func (log *CTopLogger) Exit() {
	exited = true
	wg.Wait()
}

func (log *CTopLogger) StartServer() {
	ln, err := net.Listen("unix", path)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				if err, ok := err.(net.Error); ok && err.Temporary() {
					continue
				}
				return
			}
			go log.handler(conn)
		}
	}()

	log.Notice("logging server started")
}

func (log *CTopLogger) handler(conn net.Conn) {
	wg.Add(1)
	defer wg.Done()
	defer conn.Close()
	for msg := range log.tail() {
		msg = fmt.Sprintf("%s\n", msg)
		conn.Write([]byte(msg))
	}
	conn.Write([]byte("bye\n"))
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
