package views

import (
	"strings"

	ui "github.com/gizak/termui"
)

var helpDialog = []string{
	"[h] - open this help dialog",
	"[q] - exit ctop",
}

func Help() {
	p := ui.NewPar(strings.Join(helpDialog, "\n"))
	p.Height = 10
	p.Width = 50
	p.TextFgColor = ui.ColorWhite
	p.BorderLabel = "Help"
	p.BorderFg = ui.ColorCyan
	ui.Render(p)
	ui.Handle("/sys/kbd/", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Loop()
}
