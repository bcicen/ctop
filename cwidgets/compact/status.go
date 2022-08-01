package compact

import (
	"github.com/bcicen/ctop/models"

	ui "github.com/gizak/termui"
)

// Status indicator
type Status struct {
	*ui.Block
	status []ui.Cell
	health []ui.Cell
}

func NewStatus() CompactCol {
	s := &Status{
		Block:  ui.NewBlock(),
		status: []ui.Cell{{Ch: ' '}},
		health: []ui.Cell{{Ch: ' '}},
	}
	s.Height = 1
	s.Border = false
	return s
}

func (s *Status) Buffer() ui.Buffer {
	buf := s.Block.Buffer()
	buf.Set(s.InnerX(), s.InnerY(), s.health[0])
	buf.Set(s.InnerX()+2, s.InnerY(), s.status[0])
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
	color := ui.ColorDefault
	var mark string

	switch val {
	case "":
		return
	case "created":
		mark = "◉"
	case "running":
		mark = "▶"
		color = ui.ThemeAttr("status.ok")
	case "exited":
		mark = "⏹"
		color = ui.ThemeAttr("status.danger")
	case "paused":
		mark = "⏸"
	default:
		mark = " "
		log.Warningf("unknown status string: \"%v\"", val)
	}

	s.status = ui.TextCells(mark, color, ui.ColorDefault)
}

func (s *Status) setHealth(val string) {
	color := ui.ColorDefault
	var mark string

	switch val {
	case "":
		return
	case "healthy":
		mark = "☼"
		color = ui.ThemeAttr("status.ok")
	case "unhealthy":
		mark = "⚠"
		color = ui.ThemeAttr("status.danger")
	case "starting":
		mark = "◌"
		color = ui.ThemeAttr("status.warn")
	default:
		mark = " "
		log.Warningf("unknown health state string: \"%v\"", val)
	}

	s.health = ui.TextCells(mark, color, ui.ColorDefault)
}
