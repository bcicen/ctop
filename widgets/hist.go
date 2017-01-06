package widgets

type HistData struct {
	data    []int
	labels  []string
	maxSize int
}

func NewHistData(max int) HistData {
	return HistData{
		data:    make([]int, max),
		labels:  make([]string, max),
		maxSize: max,
	}
}

func (h HistData) Append(val int) {
	if len(h.data) >= h.maxSize {
		h.data = append(h.data[:0], h.data[1:]...)
	}
	h.data = append(h.data, val)
}
