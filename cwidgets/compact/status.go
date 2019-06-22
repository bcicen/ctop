package compact

import (
	ui "github.com/gizak/termui"
)

const (
	mark       = "◉"
	healthMark = "✚"
	vBar       = string('\u25AE') + string('\u25AE')
)

// Status indicator
type Status struct {
	*ui.Block
	status []ui.Cell
	health []ui.Cell
}

func NewStatus() *Status {
	s := &Status{
		Block:  ui.NewBlock(),
		health: []ui.Cell{{Ch: ' '}},
	}
	s.Height = 1
	s.Border = false
	s.Set("")
	return s
}

func (s *Status) Buffer() ui.Buffer {
	buf := s.Block.Buffer()
	x := 0
	for _, c := range s.health {
		buf.Set(s.InnerX()+x, s.InnerY(), c)
		x += c.Width()
	}
	x += 1
	for _, c := range s.status {
		buf.Set(s.InnerX()+x, s.InnerY(), c)
		x += c.Width()
	}
	return buf
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
		text = vBar
	}

	s.status = ui.TextCells(text, color, ui.ColorDefault)
}

func (s *Status) SetHealth(val string) {
	if val == "" {
		return
	}

	color := ui.ColorDefault
	mark := healthMark

	switch val {
	case "healthy":
		color = ui.ThemeAttr("status.ok")
	case "unhealthy":
		color = ui.ThemeAttr("status.danger")
	case "starting":
		color = ui.ThemeAttr("status.warn")
	}

	s.health = ui.TextCells(mark, color, ui.ColorDefault)
}
