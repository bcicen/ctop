package compact

import (
	"github.com/bcicen/ctop/models"

	ui "github.com/gizak/termui"
)

// Status indicator
type TextStatus struct {
	*ui.Block
	status []ui.Cell
	health []ui.Cell
}

func NewTextStatus() CompactCol {
	s := &TextStatus{
		Block:  ui.NewBlock(),
		status: []ui.Cell{{Ch: ' '}},
		health: []ui.Cell{{Ch: ' '}},
	}
	s.Height = 1
	s.Border = false
	return s
}

func (s *TextStatus) Buffer() ui.Buffer {
	buf := s.Block.Buffer()
	buf.Set(s.InnerX(), s.InnerY(), s.health[0])
	buf.Set(s.InnerX()+2, s.InnerY(), s.status[0])
	return buf
}

func (s *TextStatus) SetMeta(m models.Meta) {
	s.setState(m.Get("state"))
	s.setHealth(m.Get("health"))
}

// Status implements CompactCol
func (s *TextStatus) Reset()                    {}
func (s *TextStatus) SetMetrics(models.Metrics) {}
func (s *TextStatus) Highlight()                {}
func (s *TextStatus) UnHighlight()              {}
func (s *TextStatus) Header() string            { return "" }
func (s *TextStatus) FixedWidth() int           { return 2 }

func (s *TextStatus) setState(val string) {
	color := ui.ColorDefault
	var mark string

	switch val {
	case "":
		return
	case "created":
		mark = "C"
		color = ui.ThemeAttr("fg")
	case "running":
		mark = "R"
		color = ui.ThemeAttr("status.ok")
	case "exited":
		mark = "X"
		color = ui.ThemeAttr("status.exited")
	case "paused":
		mark = "P"
		color = ui.ThemeAttr("status.paused")
	default:
		mark = " "
		log.Warningf("unknown status string: \"%v\"", val)
	}

	s.status = ui.TextCells(mark, color, ui.ColorDefault)
}

func (s *TextStatus) setHealth(val string) {
	color := ui.ColorDefault
	var mark string

	switch val {
	case "":
		return
	case "healthy":
		mark = "H"
		color = ui.ThemeAttr("status.ok")
	case "unhealthy":
		mark = "!"
		color = ui.ThemeAttr("status.danger")
	case "starting":
		mark = "s"
		color = ui.ThemeAttr("fg")
	default:
		mark = " "
		log.Warningf("unknown health state string: \"%v\"", val)
	}

	s.health = ui.TextCells(mark, color, ui.ColorDefault)
}
