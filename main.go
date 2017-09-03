package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/connector"
	"github.com/bcicen/ctop/cwidgets/compact"
	"github.com/bcicen/ctop/logging"
	"github.com/bcicen/ctop/widgets"
	ui "github.com/gizak/termui"
	tm "github.com/nsf/termbox-go"
	"github.com/bcicen/ctop/entity"
	"github.com/bcicen/ctop/network"
)

var (
	build     = "none"
	version   = "dev-build"
	goVersion = runtime.Version()

	log    *logging.CTopLogger
	cursor *GridCursor
	cGrid  *compact.CompactGrid
	header *widgets.CTopHeader

	versionStr = fmt.Sprintf("ctop version %v, build %v %v", version, build, goVersion)
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
	var invertFlag = flag.Bool("i", false, "invert default colors")
	var swarmFlag = flag.Bool("w", false, "enable s(W)arm mode")
	var imageFlag = flag.String("I", "", "name images for build service in swarm mode")
	var displayFlag = flag.Bool("H", false, "enable/disable display for service in swarm mode")
	var connectorFlag = flag.String("connector", "docker", "container connector to use")
	flag.Parse()

	if *versionFlag {
		fmt.Println(versionStr)
		os.Exit(0)
	}

	if *helpFlag {
		printHelp()
		os.Exit(0)
	}

	// init logger
	log = logging.Init()

	// init global config
	config.Init()

	// override default config values with command line flags
	if *filterFlag != "" {
		config.Update("filterStr", *filterFlag)
	}

	if *activeOnlyFlag {
		config.Toggle("allContainers")
	}

	if *sortFieldFlag != "" {
		validSort(*sortFieldFlag)
		config.Update("sortField", *sortFieldFlag)
	}

	if *reverseSortFlag {
		config.Toggle("sortReversed")
	}

	// init ui
	if *invertFlag {
		InvertColorMap()
	}

	if *swarmFlag {
		config.Toggle("swarmMode")
	}

	if *imageFlag != "" {
		config.Update("image", *imageFlag)
	}

	if *displayFlag {
		fmt.Printf("Start mode withou output.... CTRL+C for exit.")
		config.Toggle("hideDisplay")
	}

	if config.GetSwitchVal("hideDisplay") {
		ui.ColorMap = ColorMap // override default colormap
		if err := ui.Init(); err != nil {
			panic(fmt.Sprintf("Ui error: %s:", err))
		}
	} else {
		network.Main()
	}

	defer Shutdown()
	// init grid, cursor, header
	conn, err := connector.ByName(*connectorFlag)
	if err != nil {
		panic(fmt.Sprintf("Init grid, cursor, header: %s", err))
	}
	cursor = &GridCursor{cSource: conn}
	cGrid = compact.NewCompactGrid()
	header = widgets.NewCTopHeader()

	for {
		if config.GetSwitchVal("hideDisplay") {
			exit := Display()
			if exit {
				return
			}
		}
	}
}

func Shutdown() {
	if config.GetSwitchVal("swarmMode") {
		cursor.cSource.DownSwarmMode()
	}
	log.Notice("shutting down")
	log.Exit()
	if tm.IsInit {
		ui.Close()
	}
}

// ensure a given sort field is valid
func validSort(s string) {
	if _, ok := entity.Sorters[s]; !ok {
		fmt.Printf("invalid sort field: %s\n", s)
		os.Exit(1)
	}
}

func panicExit() {
	if r := recover(); r != nil {
		Shutdown()
		fmt.Printf("error: %s\n", r)
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
