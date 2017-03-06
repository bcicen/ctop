package expanded

import (
	"github.com/bcicen/ctop/metrics"
	ui "github.com/gizak/termui"
)

type Expanded struct {
	Info    *Info
	Net     *ExpandedNet
	Cpu     *ExpandedCpu
	Mem     *ExpandedMem
	infoMap map[string]string
}

func NewExpanded(id string) *Expanded {
	return &Expanded{
		Info: NewInfo(id),
		Net:  NewExpandedNet(),
		Cpu:  NewExpandedCpu(),
		Mem:  NewExpandedMem(),
	}
}

func (w *Expanded) Buffer() ui.Buffer {
	buf := ui.NewBuffer()
	buf.Merge(w.Info.Buffer())
	buf.Merge(w.Cpu.Buffer())
	buf.Merge(w.Mem.Buffer())
	buf.Merge(w.Net.Buffer())
	return buf
}

func (w *Expanded) SetMeta(k, v string) {
	w.Info.Set(k, v)
}

func (w *Expanded) SetMetrics(m metrics.Metrics) {
	w.Cpu.Update(m.CPUUtil)
	w.Net.Update(m.NetRx, m.NetTx)
	w.Mem.Update(int(m.MemUsage), int(m.MemLimit))
}
