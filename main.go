package main

import (
	"flag"
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
	defer panicExit()

	// parse command line arguments
	var versionFlag = flag.Bool("v", false, "output version information and exit")
	var helpFlag = flag.Bool("h", false, "display this help dialog")
	var filterFlag = flag.String("f", "", "filter containers")
	var activeOnlyFlag = flag.Bool("a", false, "show active containers only")
	var sortFieldFlag = flag.String("s", "", "select container sort field")
	var reverseSortFlag = flag.Bool("r", false, "reverse container sort order")
	flag.Parse()

	if *versionFlag == true {
		printVersion()
		os.Exit(0)
	}

	if *helpFlag == true {
		printHelp()
		os.Exit(0)
	}

	// init ui
	ui.ColorMap = ColorMap // override default colormap
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	// init global config
	config.Init()

	// override default config values with command line flags
	if *filterFlag != "" {
		config.Update("filterStr", *filterFlag)
	}

	if *activeOnlyFlag == true {
		config.Toggle("allContainers")
	}

	if *sortFieldFlag != "" {
		config.Update("sortField", *sortFieldFlag)
	}

	if *reverseSortFlag == true {
		config.Toggle("sortReversed")
	}

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

func panicExit() {
	if r := recover(); r != nil {
		ui.Clear()
		fmt.Printf("panic: %s\n", r)
		os.Exit(1)
	}
}

var helpMsg = `ctop - container metric viewer

usage: ctop [options]

options:
`

func printHelp() {
	fmt.Println(helpMsg)
	flag.PrintDefaults()
}

func printVersion() {
	fmt.Printf("ctop version %v, build %v\n", version, build)
}
