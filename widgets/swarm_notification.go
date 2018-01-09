package widgets

import (
	"path"

	"github.com/bcicen/ctop/config"
	ui "github.com/gizak/termui"
)

func ShowNotifiation() {
	p := ui.NewPar("You run ctop in swarm mode. \nBut agents not discovered. \nDo you want deploy agents (y/n)")
	p.Height = 5
	p.Width = 50
	p.BorderLabel = "Notification"
	ui.Render(p)
	ui.Handle("/sys/kbd/y", func(ui.Event) {
		ui.StopLoop()
		delete(ui.DefaultEvtStream.Handlers, path.Clean("/sys/kbd/y"))
	})
	ui.Handle("/sys/kbd/n", func(ui.Event) {
		ui.StopLoop()
		config.Toggle("swarmMode")
		delete(ui.DefaultEvtStream.Handlers, path.Clean("/sys/kbd/n"))
	})
	ui.Loop()
}
