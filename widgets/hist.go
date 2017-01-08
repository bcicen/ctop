package widgets

type HistData struct {
	maxLen int
	labels []string
}

func NewHistData(max int) HistData {
	return HistData{
		maxLen: max,
		labels: make([]string, max),
	}
}

type IntHistData struct {
	HistData
	data []int
}

func NewIntHistData(max int) IntHistData {
	return IntHistData{NewHistData(max), make([]int, max)}
}

func (h IntHistData) Append(val int) {
	if len(h.data) >= h.maxLen {
		h.data = append(h.data[:0], h.data[1:]...)
	}

	h.data = append(h.data, val)
}

type FloatHistData struct {
	HistData
	data []float64
}

func NewFloatHistData(max int) FloatHistData {
	return FloatHistData{NewHistData(max), make([]float64, max)}
}

func (h FloatHistData) Append(val float64) {
	if len(h.data) >= h.maxLen {
		h.data = append(h.data[:0], h.data[1:]...)
	}
	h.data = append(h.data, val)
}

type DiffHistData struct {
	HistData
	data    []int
	srcData []int
}

func NewDiffHistData(max int) DiffHistData {
	return DiffHistData{
		NewHistData(max),
		make([]int, max),
		make([]int, max),
	}
}

// return most recent value
func (h DiffHistData) Last() int {
	return h.data[len(h.data)-1]
}

func (h DiffHistData) Append(val int) {
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
