package main

import (
	"fmt"

	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/cwidgets/compact"
	"github.com/bcicen/ctop/widgets"
	ui "github.com/gizak/termui"
)

var (
	cGrid  = compact.NewCompactGrid()
	header = widgets.NewCTopHeader()
)

func maxRows() int {
	return ui.TermHeight() - 2 - cGrid.Y
}

type Grid struct {
	cursorID   string // id of currently selected container
	cSource    ContainerSource
	containers Containers // sorted slice of containers
}

func NewGrid() *Grid {
	g := &Grid{
		//cSource: NewDockerContainerSource(),
		cSource: NewMockContainerSource(),
	}
	return g
}

// Set an initial cursor position, if possible
func (g *Grid) cursorReset() {
	if len(g.containers) > 0 {
		g.cursorID = g.containers[0].Id
		g.containers[0].Widgets.Name.Highlight()
	}
}

// Return current cursor index
func (g *Grid) cursorIdx() int {
	for n, c := range g.containers {
		if c.Id == g.cursorID {
			return n
		}
	}
	return 0
}

func (g *Grid) cursorUp() {
	idx := g.cursorIdx()
	// decrement if possible
	if idx <= 0 {
		return
	}
	active := g.containers[idx]
	next := g.containers[idx-1]

	active.Widgets.Name.UnHighlight()
	g.cursorID = next.Id
	next.Widgets.Name.Highlight()

	ui.Render(cGrid)
}

func (g *Grid) cursorDown() {
	idx := g.cursorIdx()
	// increment if possible
	if idx >= (len(g.containers) - 1) {
		return
	}
	if idx >= maxRows()-1 {
		return
	}
	active := g.containers[idx]
	next := g.containers[idx+1]

	active.Widgets.Name.UnHighlight()
	g.cursorID = next.Id
	next.Widgets.Name.Highlight()
	ui.Render(cGrid)
}

func (g *Grid) redrawRows() {
	// reinit body rows
	cGrid.Clear()

	// build layout
	y := 1
	if config.GetSwitchVal("enableHeader") {
		header.SetCount(len(g.containers))
		header.SetFilter(config.GetVal("filterStr"))
		y += header.Height()
	}
	cGrid.SetY(y)

	var cursorVisible bool
	max := maxRows()
	for n, c := range g.containers {
		if n >= max {
			break
		}
		cGrid.AddRows(c.Widgets)
		if c.Id == g.cursorID {
			cursorVisible = true
		}
	}

	if !cursorVisible {
		g.cursorReset()
	}

	ui.Clear()
	if config.GetSwitchVal("enableHeader") {
		header.Render()
	}
	cGrid.Align()
	ui.Render(cGrid)
}

func (g *Grid) refreshContainers() {
	g.containers = g.cSource.All().Filter()
}

// Log current container and widget state
func (g *Grid) dumpContainer() {
	c, _ := g.cSource.Get(g.cursorID)
	msg := fmt.Sprintf("logging state for container: %s\n", c.Id)
	msg += fmt.Sprintf("Id = %s\nname = %s\nstate = %s\n", c.Id, c.Name, c.State)
	msg += inspect(&c.Metrics)
	log.Infof(msg)
}

//func (g *Grid) ExpandView() {
//ui.Clear()
//ui.DefaultEvtStream.ResetHandlers()
//defer ui.DefaultEvtStream.ResetHandlers()

//container, _ := g.cSource.Get(g.cursorID)
//// copy current widgets to restore on exit view
//curWidgets := container.widgets
//container.Expand()

//ui.Render(container.widgets)
//ui.Handle("/timer/1s", func(ui.Event) {
//ui.Render(container.widgets)
//})
//ui.Handle("/sys/kbd/", func(ui.Event) {
//ui.StopLoop()
//})
//ui.Loop()

//container.widgets = curWidgets
//container.widgets.Reset()
//}

func Display(g *Grid) bool {
	var menu func()

	cGrid.SetWidth(ui.TermWidth())
	ui.DefaultEvtStream.Hook(logEvent)

	// initial draw
	header.Align()
	g.refreshContainers()
	g.redrawRows()

	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		g.cursorUp()
	})
	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		g.cursorDown()
	})
	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		//c := g.containers[g.cursorIdx()]
		//c.Widgets.ToggleExpand()
		g.redrawRows()
	})

	ui.Handle("/sys/kbd/a", func(ui.Event) {
		config.Toggle("allContainers")
		g.refreshContainers()
		g.redrawRows()
	})
	ui.Handle("/sys/kbd/D", func(ui.Event) {
		g.dumpContainer()
	})
	ui.Handle("/sys/kbd/f", func(ui.Event) {
		menu = FilterMenu
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/h", func(ui.Event) {
		menu = HelpMenu
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/H", func(ui.Event) {
		config.Toggle("enableHeader")
		g.redrawRows()
	})
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/r", func(e ui.Event) {
		config.Toggle("sortReversed")
	})
	ui.Handle("/sys/kbd/s", func(ui.Event) {
		menu = SortMenu
		ui.StopLoop()
	})

	ui.Handle("/timer/1s", func(e ui.Event) {
		g.refreshContainers()
		g.redrawRows()
	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		header.Align()
		cGrid.SetWidth(ui.TermWidth())
		log.Infof("resize: width=%v max-rows=%v", cGrid.Width, maxRows())
		g.redrawRows()
	})

	ui.Loop()
	if menu != nil {
		ui.Clear()
		menu()
		return false
	}
	return true
}
