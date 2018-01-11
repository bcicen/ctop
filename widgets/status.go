package widgets

import (
	ui "github.com/gizak/termui"
)

var (
	statusHeight = 1
	statusIter   = 3
)

type StatusLine struct {
	Message *ui.Par
	bg      *ui.Par
}

func NewStatusLine() *StatusLine {
	p := ui.NewPar("")
	p.X = 2
	p.Border = false
	p.Height = statusHeight
	p.Bg = ui.ThemeAttr("header.bg")
	p.TextFgColor = ui.ThemeAttr("header.fg")
	p.TextBgColor = ui.ThemeAttr("header.bg")
	return &StatusLine{
		Message: p,
		bg:      statusBg(),
	}
}

func (sl *StatusLine) Display() {
	ui.DefaultEvtStream.ResetHandlers()
	defer ui.DefaultEvtStream.ResetHandlers()

	iter := statusIter
	ui.Handle("/sys/kbd/", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/timer/1s", func(ui.Event) {
		iter--
		if iter <= 0 {
			ui.StopLoop()
		}
	})

	ui.Render(sl)
	ui.Loop()
}

// change given message on the status line
func (sl *StatusLine) Show(s string) {
	sl.Message.TextFgColor = ui.ThemeAttr("header.fg")
	sl.Message.Text = s
	sl.Display()
}

func (sl *StatusLine) ShowErr(s string) {
	sl.Message.TextFgColor = ui.ThemeAttr("status.danger")
	sl.Message.Text = s
	sl.Display()
}

func (sl *StatusLine) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	buf.Merge(sl.bg.Buffer())
	buf.Merge(sl.Message.Buffer())
	return buf
}

func (sl *StatusLine) Align() {
	sl.bg.SetWidth(ui.TermWidth() - 1)
	sl.Message.SetWidth(ui.TermWidth() - 2)

	sl.bg.Y = ui.TermHeight() - 1
	sl.Message.Y = ui.TermHeight() - 1
}

func (sl *StatusLine) Height() int { return statusHeight }

func statusBg() *ui.Par {
	bg := ui.NewPar("")
	bg.X = 1
	bg.Height = statusHeight
	bg.Border = false
	bg.Bg = ui.ThemeAttr("header.bg")
	return bg
}
