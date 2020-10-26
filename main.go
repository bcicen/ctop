package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/connector"
	"github.com/bcicen/ctop/container"
	"github.com/bcicen/ctop/cwidgets/compact"
	"github.com/bcicen/ctop/logging"
	"github.com/bcicen/ctop/widgets"
	ui "github.com/gizak/termui"
	tm "github.com/nsf/termbox-go"
)

var (
	build     = "none"
	version   = "dev-build"
	goVersion = runtime.Version()

	log     *logging.CTopLogger
	cursor  *GridCursor
	cGrid   *compact.CompactGrid
	header  *widgets.CTopHeader
	status  *widgets.StatusLine
	errView *widgets.ErrorView

	versionStr = fmt.Sprintf("ctop version %v, build %v %v", version, build, goVersion)
)

func main() {
	defer panicExit()

	// parse command line arguments
	var (
		versionFlag     = flag.Bool("v", false, "output version information and exit")
		helpFlag        = flag.Bool("h", false, "display this help dialog")
		filterFlag      = flag.String("f", "", "filter containers")
		activeOnlyFlag  = flag.Bool("a", false, "show active containers only")
		sortFieldFlag   = flag.String("s", "", "select container sort field")
		reverseSortFlag = flag.Bool("r", false, "reverse container sort order")
		invertFlag      = flag.Bool("i", false, "invert default colors")
		scaleCpu        = flag.Bool("scale-cpu", false, "show cpu as % of system total")
		connectorFlag   = flag.String("connector", "docker", "container connector to use")
		defaultShell    = flag.String("shell", "sh", "exec shell to use")
	)
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

	// init global config and read config file if exists
	config.Init()
	config.Read()

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

	if *scaleCpu {
		config.Toggle("scaleCpu")
	}

	if *defaultShell != "" {
		config.Update("shell", *defaultShell)
	}

	// init ui
	if *invertFlag {
		InvertColorMap()
	}
	ui.ColorMap = ColorMap // override default colormap
	if err := ui.Init(); err != nil {
		panic(err)
	}
	tm.SetInputMode(tm.InputAlt)

	defer Shutdown()
	// init grid, cursor, header
	cSuper, err := connector.ByName(*connectorFlag)
	if err != nil {
		panic(err)
	}
	cursor = &GridCursor{cSuper: cSuper}
	cGrid = compact.NewCompactGrid()
	header = widgets.NewCTopHeader()
	status = widgets.NewStatusLine()
	errView = widgets.NewErrorView()

	for {
		exit := Display()
		if exit {
			return
		}
	}
}

func Shutdown() {
	log.Notice("shutting down")
	log.Exit()
	if tm.IsInit {
		ui.Close()
	}
}

// ensure a given sort field is valid
func validSort(s string) {
	if _, ok := container.Sorters[s]; !ok {
		fmt.Printf("invalid sort field: %s\n", s)
		os.Exit(1)
	}
}

func panicExit() {
	if r := recover(); r != nil {
		Shutdown()
		panic(r)
		fmt.Printf("error: %s\n", r)
		os.Exit(1)
	}
}

var helpMsg = `ctop - interactive container viewer

usage: ctop [options]

options:
`

func printHelp() {
	fmt.Println(helpMsg)
	flag.PrintDefaults()
	fmt.Printf("\navailable connectors: ")
	fmt.Println(strings.Join(connector.Enabled(), ", "))
}
