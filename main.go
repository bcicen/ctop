package main

import (
	"fmt"
	"os"

	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/logging"
	ui "github.com/gizak/termui"
)

var log *logging.CTopLogger

func main() {
	defer func() {
		if r := recover(); r != nil {
			ui.Clear()
			fmt.Printf("panic: %s", r)
			os.Exit(1)
		}
	}()
	config.Init()
	log = logging.Init()
	if config.GetSwitchVal("loggingEnabled") {
		logging.StartServer()
	}
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	g := NewGrid()
	for {
		exit := Display(g)
		if exit {
			log.Notice("shutting down")
			log.Exit()
			return
		}
	}
}
