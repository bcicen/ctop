package main

import (
	"github.com/bcicen/ctop/logging"
	ui "github.com/gizak/termui"
)

var log = logging.New(GlobalConfig["loggingEnabled"])

func main() {
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	g := NewGrid()
	for {
		exit := Display(g)
		if exit {
			return
		}
	}
}
