package cwidgets

import (
	"github.com/bcicen/ctop/logging"
	ui "github.com/gizak/termui"
)

var log = logging.Init()

type ContainerWidgets interface {
	Render(int, int)
	Reset()
	Buffer() ui.Buffer
	Highlight()
	UnHighlight()
	SetY(int)
	SetWidth(int)
	SetStatus(string)
	SetCPU(int)
	SetNet(int64, int64)
	SetMem(int64, int64, int)
}
