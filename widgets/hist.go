package widgets

type Hist struct {
	maxLen int
	labels []string
}

func NewHist(max int) Hist {
	return Hist{
		maxLen: max,
		labels: make([]string, max),
	}
}

type IntHist struct {
	Hist
	data []int
}

func NewIntHist(max int) IntHist {
	return IntHist{NewHist(max), make([]int, max)}
}

func (h IntHist) Append(val int) {
	if len(h.data) >= h.maxLen {
		h.data = append(h.data[:0], h.data[1:]...)
	}

	h.data = append(h.data, val)
}

type FloatHist struct {
	Hist
	data []float64
}

func NewFloatHist(max int) FloatHist {
	return FloatHist{NewHist(max), make([]float64, max)}
}

func (h FloatHist) Append(val float64) {
	if len(h.data) >= h.maxLen {
		h.data = append(h.data[:0], h.data[1:]...)
	}
	h.data = append(h.data, val)
}

type DiffHist struct {
	Hist
	data    []int
	srcData []int
}

func NewDiffHist(max int) DiffHist {
	return DiffHist{
		NewHist(max),
		make([]int, max),
		make([]int, max),
	}
}

// return most recent value
func (h DiffHist) Last() int {
	return h.data[len(h.data)-1]
}

func (h DiffHist) Append(val int) {
	if len(h.data) >= h.maxLen {
		h.data = append(h.data[:0], h.data[1:]...)
	}
	if len(h.srcData) >= h.maxLen {
		h.srcData = append(h.srcData[:0], h.srcData[1:]...)
	}

	diff := val - h.srcData[len(h.srcData)-1]
	if diff != val { // skip adding to data if this is the initial update
		h.data = append(h.data, diff)
	}
	h.srcData = append(h.srcData, val)
}
