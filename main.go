package main

import (
	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/logging"
	ui "github.com/gizak/termui"
)

var log *logging.CTopLogger

func main() {
	log = logging.Init(config.Global["loggingEnabled"])
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
