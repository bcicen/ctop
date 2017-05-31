package logging

import (
	"fmt"
	"net"
	"sync"
)

const (
	path = "./ctop.sock"
)

var server struct {
	wg sync.WaitGroup
	ln net.Listener
}

func getListener() net.Listener {
	ln, err := net.Listen("unix", path)
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

func handler(conn net.Conn) {
	server.wg.Add(1)
	defer server.wg.Done()
	defer conn.Close()
	for msg := range Log.tail() {
		msg = fmt.Sprintf("%s\n", msg)
		conn.Write([]byte(msg))
	}
	conn.Write([]byte("bye\n"))
}
