package compact

import (
	"github.com/bcicen/ctop/models"
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
	s.setState("")
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

func (s *Status) SetMeta(m models.Meta) {
	s.setState(m.Get("state"))
	s.setHealth(m.Get("health"))
}

// Status implements CompactCol
func (s *Status) Reset()                    {}
func (s *Status) SetMetrics(models.Metrics) {}
func (s *Status) Highlight()                {}
func (s *Status) UnHighlight()              {}
func (s *Status) Header() string            { return "" }
func (s *Status) FixedWidth() int           { return 3 }

func (s *Status) setState(val string) {
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

func (s *Status) setHealth(val string) {
	color := ui.ColorDefault
	mark := healthMark

	switch val {
	case "":
		return
	case "healthy":
		color = ui.ThemeAttr("status.ok")
	case "unhealthy":
		color = ui.ThemeAttr("status.danger")
	case "starting":
		color = ui.ThemeAttr("status.warn")
	default:
		log.Warningf("unknown health state string: \"%v\"", val)
	}

	s.health = ui.TextCells(mark, color, ui.ColorDefault)
}
