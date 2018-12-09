package compact

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bcicen/ctop/cwidgets"
	ui "github.com/gizak/termui"
)

func (row *Compact) SetNet(rx int64, tx int64) {
	label := fmt.Sprintf("%s / %s", cwidgets.ByteFormat(rx), cwidgets.ByteFormat(tx))
	row.Net.Set(label)
}

func (row *Compact) SetIO(read int64, write int64) {
	label := fmt.Sprintf("%s / %s", cwidgets.ByteFormat(read), cwidgets.ByteFormat(write))
	row.IO.Set(label)
}

func (row *Compact) SetPids(val int) {
	label := strconv.Itoa(val)
	row.Pids.Set(label)
}

func (row *Compact) SetUptime(val int64) {
	d := time.Duration(val) * time.Millisecond
	label := "- h"
	if d.Hours() < 1.0 {
		label = fmt.Sprintf("%.0fm", d.Minutes())
	} else if d.Hours() < 24.0 {
		label = fmt.Sprintf("%.0fh", d.Hours())
	} else {
		label = fmt.Sprintf("%dd", int(d.Hours())%24)
	}
	row.Uptime.Set(label)
}

func (row *Compact) SetCPU(val int) {
	row.Cpu.BarColor = colorScale(val)
	row.Cpu.Label = fmt.Sprintf("%s%%", strconv.Itoa(val))
	if val < 5 {
		val = 5
		row.Cpu.BarColor = ui.ThemeAttr("gauge.bar.bg")
	}
	if val > 100 {
		val = 100
	}
	row.Cpu.Percent = val
}

func (row *Compact) SetMem(val int64, limit int64, percent int) {
	row.Mem.Label = fmt.Sprintf("%s / %s", cwidgets.ByteFormat(val), cwidgets.ByteFormat(limit))
	if percent < 5 {
		percent = 5
		row.Mem.BarColor = ui.ColorBlack
	} else {
		row.Mem.BarColor = ui.ThemeAttr("gauge.bar.bg")
	}
	row.Mem.Percent = percent
}
