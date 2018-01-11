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
	{"[q] - exit ctop", ""},
}

func HelpMenu() {
	ui.Clear()
	ui.DefaultEvtStream.ResetHandlers()
	defer ui.DefaultEvtStream.ResetHandlers()

	m := menu.NewMenu()
	m.BorderLabel = "Help"
	m.AddItems(helpDialog...)
	ui.Render(m)
	ui.Handle("/sys/kbd/", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Loop()
}

func FilterMenu() {
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
}

func SortMenu() {
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
}

func ContainerMenu() {

	c := cursor.Selected()
	if c == nil {
		return
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
	if c.Meta["state"] == "exited" {
		items = append(items, menu.Item{Val: "start", Label: "start"})
		items = append(items, menu.Item{Val: "remove", Label: "remove"})
	}
	items = append(items, menu.Item{Val: "cancel", Label: "cancel"})

	m.AddItems(items...)
	ui.Render(m)

	confirmTxt := func(a, n string) string { return fmt.Sprintf("%s container %s?", a, n) }

	HandleKeys("up", m.Up)
	HandleKeys("down", m.Down)
	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		switch m.SelectedItem().Val {
		case "single":
			SingleView(c)
			ui.StopLoop()
		case "logs":
			LogMenu()
			ui.StopLoop()
		case "start":
			Confirm(confirmTxt("start", c.GetMeta("name")), c.Start)
			ui.StopLoop()
		case "stop":
			Confirm(confirmTxt("stop", c.GetMeta("name")), c.Stop)
			ui.StopLoop()
		case "remove":
			Confirm(confirmTxt("remove", c.GetMeta("name")), c.Remove)
			ui.StopLoop()
		case "cancel":
			ui.StopLoop()
		}
	})
	ui.Handle("/sys/kbd/", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Loop()
}

func LogMenu() {

	c := cursor.Selected()
	if c == nil {
		return
	}

	ui.DefaultEvtStream.ResetHandlers()
	defer ui.DefaultEvtStream.ResetHandlers()

	logs, quit := logReader(c)
	m := widgets.NewTextView(logs)
	m.BorderLabel = "Logs"
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
}

func Confirm(txt string, fn func()) {
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
