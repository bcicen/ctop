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
		config.Update("sortField", m.SelectedItem().Val)
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
		menu.Item{Val: "single", Label: "single view"},
		menu.Item{Val: "logs", Label: "log view"},
	}

	if c.Meta["state"] == "running" {
		items = append(items, menu.Item{Val: "stop", Label: "stop"})
	}
	if c.Meta["state"] == "exited" || c.Meta["state"] == "created" {
		items = append(items, menu.Item{Val: "start", Label: "start"})
		items = append(items, menu.Item{Val: "remove", Label: "remove"})
	}
	items = append(items, menu.Item{Val: "cancel", Label: "cancel"})

	m.AddItems(items...)
	ui.Render(m)

	var nextMenu MenuFn
	HandleKeys("up", m.Up)
	HandleKeys("down", m.Down)
	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		switch m.SelectedItem().Val {
		case "single":
			nextMenu = SingleView
		case "logs":
			nextMenu = LogMenu
		case "start":
			nextMenu = Confirm(confirmTxt("start", c.GetMeta("name")), c.Start)
		case "stop":
			nextMenu = Confirm(confirmTxt("stop", c.GetMeta("name")), c.Stop)
		case "remove":
			nextMenu = Confirm(confirmTxt("remove", c.GetMeta("name")), c.Remove)
		}
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Loop()
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
			switch m.SelectedItem().Val {
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
