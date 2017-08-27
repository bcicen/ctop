package compact

import (
	"fmt"

	ui "github.com/gizak/termui"
)

const (
	mark        = string('\u25C9')
	vBar        = string('\u25AE')
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
	case "running":
		color = ui.ThemeAttr("status.ok")
	case "exited":
		color = ui.ThemeAttr("status.danger")
	case "paused":
		text = fmt.Sprintf("%s%s", vBar, vBar)
	}

	s.Text = text
	s.TextFgColor = color
}
