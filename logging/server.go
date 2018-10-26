package logging

import (
	"fmt"
	"io"
	"net"
	"sync"
)

const (
	socketPath = "./ctop.sock"
	socketAddr = "0.0.0.0:9000"
)

var server struct {
	wg sync.WaitGroup
	ln net.Listener
}

func getListener() net.Listener {
	var ln net.Listener
	var err error
	if debugModeTCP() {
		ln, err = net.Listen("tcp", socketAddr)
	} else {
		ln, err = net.Listen("unix", socketPath)
	}
	if err != nil {
		panic(err)
	}
	return ln
}

func StartServer() {
	server.ln = getListener()

	go func() {
		for {
			conn, err := server.ln.Accept()
			if err != nil {
				if err, ok := err.(net.Error); ok && err.Temporary() {
					continue
				}
				return
			}
			go handler(conn)
		}
	}()

	Log.Notice("logging server started")
}

func StopServer() {
	server.wg.Wait()
	if server.ln != nil {
		server.ln.Close()
	}
}

func handler(wc io.WriteCloser) {
	server.wg.Add(1)
	defer server.wg.Done()
	defer wc.Close()
	for msg := range Log.tail() {
		msg = fmt.Sprintf("%s\n", msg)
		wc.Write([]byte(msg))
	}
	wc.Write([]byte("bye\n"))
}
