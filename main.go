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
	build   = "none"
	version = "dev-build"

	log    *logging.CTopLogger
	cursor *GridCursor
	cGrid  *compact.CompactGrid
	header *widgets.CTopHeader
)

func main() {
	readArgs()
	defer panicExit()

	// init ui
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	// init global config
	config.Init()

	// init logger
	log = logging.Init()
	if config.GetSwitchVal("loggingEnabled") {
		logging.StartServer()
	}

	// init grid, cursor, header
	cursor = NewGridCursor()
	cGrid = compact.NewCompactGrid()
	header = widgets.NewCTopHeader()

	for {
		exit := Display()
		if exit {
			log.Notice("shutting down")
			log.Exit()
			return
		}
	}
}

func readArgs() {
	if len(os.Args) < 2 {
		return
	}
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-v", "version":
			printVersion()
			os.Exit(0)
		case "-h", "help":
			printHelp()
			os.Exit(0)
		default:
			fmt.Printf("invalid option or argument: \"%s\"\n", arg)
			os.Exit(1)
		}
	}
}

func panicExit() {
	if r := recover(); r != nil {
		ui.Clear()
		fmt.Printf("panic: %s\n", r)
		os.Exit(1)
	}
}

var helpMsg = `cTop - container metric viewer

usage: ctop [options]

options:
 -h display this help dialog
 -v output version information and exit
`

func printHelp() {
	fmt.Println(helpMsg)
}

func printVersion() {
	fmt.Printf("cTop version %v, build %v\n", version, build)
}
