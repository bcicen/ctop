package compact

import (
	"fmt"

	ui "github.com/gizak/termui"
)

const (
	mark        = string('\u25C9')
	vBar        = string('\u25AE')
	service     = string('\u0053')
	statusWidth = 3
)

// Status indicator
type Status struct {
	*ui.Par
}

func NewStatus() *Status {
	p := ui.NewPar(mark)
	p.Border = false
	p.Height = 1
	p.Width = statusWidth
	return &Status{p}
}

func (s *Status) Set(val string) {
	// defaults
	text := mark
	color := ui.ColorDefault

	switch val {
	case "new":
		color = ui.ColorCyan
	case "running", "rollback_completed":
		color = ui.ColorGreen
	case "rollback_started", "rollback_paused", "updating", "starting", "ready":
		color = ui.ColorYellow
	case "exited", "shutdown":
		color = ui.ColorRed
	case "failed":
		color = ui.ColorMagenta
	case "paused":
		text = fmt.Sprintf("%s%s", vBar, vBar)
	case "service":
		text = fmt.Sprintf("%s", service)
	}

	s.Text = text
	s.TextFgColor = color
}
