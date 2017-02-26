package expanded

type IntHist struct {
	data   []int
	labels []string
}

func NewIntHist(max int) IntHist {
	return IntHist{
		data:   make([]int, max),
		labels: make([]string, max),
	}
}

func (h IntHist) Append(val int) {
	if len(h.data) == cap(h.data) {
		h.data = append(h.data[:0], h.data[1:]...)
	}

	h.data = append(h.data, val)
}

type FloatHist struct {
	data   []float64
	labels []string
}

func NewFloatHist(max int) FloatHist {
	return FloatHist{
		data:   make([]float64, max),
		labels: make([]string, max),
	}
}

func (h FloatHist) Append(val float64) {
	if len(h.data) == cap(h.data) {
		h.data = append(h.data[:0], h.data[1:]...)
	}
	h.data = append(h.data, val)
}

type DiffHist struct {
	data    []int // data point derivatives
	srcData []int // principal input data
	labels  []string
}

func NewDiffHist(max int) DiffHist {
	return DiffHist{
		data:    make([]int, max),
		srcData: make([]int, max),
		labels:  make([]string, max),
	}
}

// return most recent value
func (h DiffHist) Last() int {
	return h.data[len(h.data)-1]
}

func (h DiffHist) Append(val int) {
	if len(h.data) == cap(h.data) {
		h.data = append(h.data[:0], h.data[1:]...)
	}
	if len(h.srcData) == cap(h.srcData) {
		h.srcData = append(h.srcData[:0], h.srcData[1:]...)
	}

	diff := val - h.srcData[len(h.srcData)-1]
	if diff != val { // skip adding to data if this is the initial update
		h.data = append(h.data, diff)
	}
	h.srcData = append(h.srcData, val)
}
