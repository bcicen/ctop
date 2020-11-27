package logging

import (
	"fmt"
	"os"
	"time"

	"github.com/op/go-logging"
)

const (
	size = 1024
)

var (
	Log    *CTopLogger
	exited bool
	level  = logging.INFO // default level
	format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
)

type statusMsg struct {
	Text    string
	IsError bool
}

type CTopLogger struct {
	*logging.Logger
	backend *logging.MemoryBackend
	logFile *os.File
	sLog    []statusMsg
}

func (c *CTopLogger) FlushStatus() chan statusMsg {
	ch := make(chan statusMsg)
	go func() {
		for _, sm := range c.sLog {
			ch <- sm
		}
		close(ch)
		c.sLog = []statusMsg{}
	}()
	return ch
}

func (c *CTopLogger) StatusQueued() bool     { return len(c.sLog) > 0 }
func (c *CTopLogger) Status(s string)        { c.addStatus(statusMsg{s, false}) }
func (c *CTopLogger) StatusErr(err error)    { c.addStatus(statusMsg{err.Error(), true}) }
func (c *CTopLogger) addStatus(sm statusMsg) { c.sLog = append(c.sLog, sm) }

func (c *CTopLogger) Statusf(s string, a ...interface{}) { c.Status(fmt.Sprintf(s, a...)) }

func Init() *CTopLogger {
	if Log == nil {
		logging.SetFormatter(format) // setup default formatter

		Log = &CTopLogger{
			logging.MustGetLogger("ctop"),
			logging.NewMemoryBackend(size),
			nil,
			[]statusMsg{},
		}

		debugMode := debugMode()
		if debugMode {
			level = logging.DEBUG
		}
		backendLvl := logging.AddModuleLevel(Log.backend)
		backendLvl.SetLevel(level, "")

		logFilePath := debugModeFile()
		if logFilePath == "" {
			logging.SetBackend(backendLvl)
		} else {
			logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				logging.SetBackend(backendLvl)
				Log.Error("Unable to create log file: %s", err.Error())
			} else {
				backendFile := logging.NewLogBackend(logFile, "", 0)
				backendFileLvl := logging.AddModuleLevel(backendFile)
				backendFileLvl.SetLevel(level, "")
				logging.SetBackend(backendLvl, backendFileLvl)
				Log.logFile = logFile
			}
		}

		if debugMode {
			StartServer()
		}
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
	if log.logFile != nil {
		_ = log.logFile.Close()
	}
	StopServer()
}

func debugMode() bool       { return os.Getenv("CTOP_DEBUG") == "1" }
func debugModeTCP() bool    { return os.Getenv("CTOP_DEBUG_TCP") == "1" }
func debugModeFile() string { return os.Getenv("CTOP_DEBUG_FILE") }
