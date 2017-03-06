package main

import (
	"fmt"
	"os"

	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/cwidgets/compact"
	"github.com/bcicen/ctop/logging"
	"github.com/bcicen/ctop/widgets"
	ui "github.com/gizak/termui"
)

var (
	log    *logging.CTopLogger
	cursor *GridCursor
	cGrid  = compact.NewCompactGrid()
	header = widgets.NewCTopHeader()
)

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
	cursor = NewGridCursor()
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	for {
		exit := Display()
		if exit {
			log.Notice("shutting down")
			log.Exit()
			return
		}
	}
}
