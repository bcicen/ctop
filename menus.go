package main

import (
	"fmt"
	"time"

	"github.com/bcicen/ctop/config"
	"github.com/bcicen/ctop/container"
	"github.com/bcicen/ctop/widgets"
	"github.com/bcicen/ctop/widgets/menu"
	ui "github.com/gizak/termui"
)

// MenuFn executes a menu window, returning the next menu or nil
type MenuFn func() MenuFn

var helpDialog = []menu.Item{
	{"<enter> - open container menu", ""},
	{"", ""},
	{"[a] - toggle display of all containers", ""},
	{"[f] - filter displayed containers", ""},
	{"[h] - open this help dialog", ""},
	{"[H] - toggle ctop header", ""},
	{"[s] - select container sort field", ""},
	{"[r] - reverse container sort order", ""},
	{"[o] - open single view", ""},
	{"[l] - view container logs ([t] to toggle timestamp when open)", ""},
	{"[e] - exec shell", ""},
	{"[S] - save current configuration to file", ""},
	{"[q] - exit ctop", ""},
}

func HelpMenu() MenuFn {
	ui.Clear()
	ui.DefaultEvtStream.ResetHandlers()
	defer ui.DefaultEvtStream.ResetHandlers()

	m := menu.NewMenu()
	m.BorderLabel = "Help"
	m.AddItems(helpDialog...)
	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Clear()
		ui.Render(m)
	})
	ui.Handle("/sys/kbd/", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Loop()
	return nil
}

func FilterMenu() MenuFn {
	ui.DefaultEvtStream.ResetHandlers()
	defer ui.DefaultEvtStream.ResetHandlers()

	i := widgets.NewInput()
	i.BorderLabel = "Filter"
	i.SetY(ui.TermHeight() - i.Height)
	i.Data = config.GetVal("filterStr")
	ui.Render(i)

	// refresh container rows on input
	stream := i.Stream()
	go func() {
		for s := range stream {
			config.Update("filterStr", s)
			RefreshDisplay()
			ui.Render(i)
		}
	}()

	i.InputHandlers()
	ui.Handle("/sys/kbd/<escape>", func(ui.Event) {
		config.Update("filterStr", "")
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		config.Update("filterStr", i.Data)
		ui.StopLoop()
	})
	ui.Loop()
	return nil
}

func SortMenu() MenuFn {
	ui.Clear()
	ui.DefaultEvtStream.ResetHandlers()
	defer ui.DefaultEvtStream.ResetHandlers()

	m := menu.NewMenu()
	m.Selectable = true
	m.SortItems = true
	m.BorderLabel = "Sort Field"

	for _, field := range container.SortFields() {
		m.AddItems(menu.Item{field, ""})
	}

	// set cursor position to current sort field
	m.SetCursor(config.GetVal("sortField"))

	HandleKeys("up", m.Up)
	HandleKeys("down", m.Down)
	HandleKeys("exit", ui.StopLoop)

	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		config.Update("sortField", m.SelectedValue())
		ui.StopLoop()
	})

	ui.Render(m)
	ui.Loop()
	return nil
}

func ColumnsMenu() MenuFn {
	ui.Clear()
	ui.DefaultEvtStream.ResetHandlers()
	defer ui.DefaultEvtStream.ResetHandlers()

	m := menu.NewMenu()
	m.Selectable = true
	m.SortItems = false
	m.BorderLabel = "Columns"

	rebuild := func() {
		m.ClearItems()
		for _, col := range config.GlobalColumns {
			txt := fmt.Sprintf("%s [%t]", col.Label, col.Enabled)
			m.AddItems(menu.Item{col.Name, txt})
		}
	}

	upFn := func() {
		config.ColumnLeft(m.SelectedValue())
		m.Up()
		rebuild()
	}

	downFn := func() {
		config.ColumnRight(m.SelectedValue())
		m.Down()
		rebuild()
	}

	toggleFn := func() {
		config.ColumnToggle(m.SelectedValue())
		rebuild()
	}

	rebuild()

	HandleKeys("up", m.Up)
	HandleKeys("down", m.Down)
	HandleKeys("enter", toggleFn)
	HandleKeys("pgup", upFn)
	HandleKeys("pgdown", downFn)

	ui.Handle("/sys/kbd/x", func(ui.Event) { toggleFn() })
	ui.Handle("/sys/kbd/<enter>", func(ui.Event) { toggleFn() })

	HandleKeys("exit", func() {
		cSource, err := cursor.cSuper.Get()
		if err == nil {
			for _, c := range cSource.All() {
				c.RecreateWidgets()
			}
		}
		ui.StopLoop()
	})

	ui.Render(m)
	ui.Loop()
	return nil
}

func ContainerMenu() MenuFn {
	c := cursor.Selected()
	if c == nil {
		return nil
	}

	ui.DefaultEvtStream.ResetHandlers()
	defer ui.DefaultEvtStream.ResetHandlers()

	m := menu.NewMenu()
	m.Selectable = true
	m.BorderLabel = "Menu"

	items := []menu.Item{
		menu.Item{Val: "single", Label: "[o] single view"},
		menu.Item{Val: "logs", Label: "[l] log view"},
	}

	if c.Meta["state"] == "running" {
		items = append(items, menu.Item{Val: "stop", Label: "[s] stop"})
		items = append(items, menu.Item{Val: "pause", Label: "[p] pause"})
		items = append(items, menu.Item{Val: "restart", Label: "[r] restart"})
		items = append(items, menu.Item{Val: "exec", Label: "[e] exec shell"})
	}
	if c.Meta["state"] == "exited" || c.Meta["state"] == "created" {
		items = append(items, menu.Item{Val: "start", Label: "[s] start"})
		items = append(items, menu.Item{Val: "remove", Label: "[R] remove"})
	}
	if c.Meta["state"] == "paused" {
		items = append(items, menu.Item{Val: "unpause", Label: "[p] unpause"})
	}
	items = append(items, menu.Item{Val: "cancel", Label: "[c] cancel"})

	m.AddItems(items...)
	ui.Render(m)

	HandleKeys("up", m.Up)
	HandleKeys("down", m.Down)

	var selected string

	// shortcuts
	ui.Handle("/sys/kbd/o", func(ui.Event) {
		selected = "single"
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/l", func(ui.Event) {
		selected = "logs"
		ui.StopLoop()
	})
	if c.Meta["state"] != "paused" {
		ui.Handle("/sys/kbd/s", func(ui.Event) {
			if c.Meta["state"] == "running" {
				selected = "stop"
			} else {
				selected = "start"
			}
			ui.StopLoop()
		})
	}
	if c.Meta["state"] != "exited" || c.Meta["state"] != "created" {
		ui.Handle("/sys/kbd/p", func(ui.Event) {
			if c.Meta["state"] == "paused" {
				selected = "unpause"
			} else {
				selected = "pause"
			}
			ui.StopLoop()
		})
	}
	if c.Meta["state"] == "running" {
		ui.Handle("/sys/kbd/e", func(ui.Event) {
			selected = "exec"
			ui.StopLoop()
		})
		ui.Handle("/sys/kbd/r", func(ui.Event) {
			selected = "restart"
			ui.StopLoop()
		})
	}
	ui.Handle("/sys/kbd/R", func(ui.Event) {
		selected = "remove"
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/c", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		selected = m.SelectedValue()
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Loop()

	var nextMenu MenuFn
	switch selected {
	case "single":
		nextMenu = SingleView
	case "logs":
		nextMenu = LogMenu
	case "exec":
		nextMenu = ExecShell
	case "start":
		nextMenu = Confirm(confirmTxt("start", c.GetMeta("name")), c.Start)
	case "stop":
		nextMenu = Confirm(confirmTxt("stop", c.GetMeta("name")), c.Stop)
	case "remove":
		nextMenu = Confirm(confirmTxt("remove", c.GetMeta("name")), c.Remove)
	case "pause":
		nextMenu = Confirm(confirmTxt("pause", c.GetMeta("name")), c.Pause)
	case "unpause":
		nextMenu = Confirm(confirmTxt("unpause", c.GetMeta("name")), c.Unpause)
	case "restart":
		nextMenu = Confirm(confirmTxt("restart", c.GetMeta("name")), c.Restart)
	}

	return nextMenu
}

func LogMenu() MenuFn {

	c := cursor.Selected()
	if c == nil {
		return nil
	}

	ui.DefaultEvtStream.ResetHandlers()
	defer ui.DefaultEvtStream.ResetHandlers()

	logs, quit := logReader(c)
	m := widgets.NewTextView(logs)
	m.BorderLabel = fmt.Sprintf("Logs [%s]", c.GetMeta("name"))
	ui.Render(m)

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		m.Resize()
	})
	ui.Handle("/sys/kbd/t", func(ui.Event) {
		m.Toggle()
	})
	ui.Handle("/sys/kbd/", func(ui.Event) {
		quit <- true
		ui.StopLoop()
	})
	ui.Loop()
	return nil
}

func ExecShell() MenuFn {
	c := cursor.Selected()

	if c == nil {
		return nil
	}

	ui.DefaultEvtStream.ResetHandlers()
	defer ui.DefaultEvtStream.ResetHandlers()

	shell := config.Get("shell")
	if err := c.Exec([]string{shell.Val, "-c", "printf '\\e[0m\\e[?25h' && clear && " + shell.Val}); err != nil {
		log.Fatal(err)
	}

	return nil
}

// Create a confirmation dialog with a given description string and
// func to perform if confirmed
func Confirm(txt string, fn func()) MenuFn {
	menu := func() MenuFn {
		ui.DefaultEvtStream.ResetHandlers()
		defer ui.DefaultEvtStream.ResetHandlers()

		m := menu.NewMenu()
		m.Selectable = true
		m.BorderLabel = "Confirm"
		m.SubText = txt

		items := []menu.Item{
			menu.Item{Val: "cancel", Label: "[c]ancel"},
			menu.Item{Val: "yes", Label: "[y]es"},
		}

		var response bool

		m.AddItems(items...)
		ui.Render(m)

		yes := func() {
			response = true
			ui.StopLoop()
		}

		no := func() {
			response = false
			ui.StopLoop()
		}

		HandleKeys("up", m.Up)
		HandleKeys("down", m.Down)
		HandleKeys("exit", no)
		ui.Handle("/sys/kbd/c", func(ui.Event) { no() })
		ui.Handle("/sys/kbd/y", func(ui.Event) { yes() })

		ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
			switch m.SelectedValue() {
			case "cancel":
				no()
			case "yes":
				yes()
			}
		})

		ui.Loop()
		if response {
			fn()
		}
		return nil
	}
	return menu
}

type toggleLog struct {
	timestamp time.Time
	message   string
}

func (t *toggleLog) Toggle(on bool) string {
	if on {
		return fmt.Sprintf("%s %s", t.timestamp.Format("2006-01-02T15:04:05.999Z07:00"), t.message)
	}
	return t.message
}

func logReader(container *container.Container) (logs chan widgets.ToggleText, quit chan bool) {

	logCollector := container.Logs()
	stream := logCollector.Stream()
	logs = make(chan widgets.ToggleText)
	quit = make(chan bool)

	go func() {
		for {
			select {
			case log := <-stream:
				logs <- &toggleLog{timestamp: log.Timestamp, message: log.Message}
			case <-quit:
				logCollector.Stop()
				close(logs)
				return
			}
		}
	}()
	return
}

func confirmTxt(a, n string) string { return fmt.Sprintf("%s container %s?", a, n) }
