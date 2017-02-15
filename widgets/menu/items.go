package menu

type Item struct {
	Val   string
	Label string
}

// Use label as display text of item, if given
func (m Item) Text() string {
	if m.Label != "" {
		return m.Label
	}
	return m.Val
}

type Items []Item

func NewItems(items ...Item) (mitems Items) {
	for _, i := range items {
		mitems = append(mitems, i)
	}
	return mitems
}

// Sort methods for Items
func (m Items) Len() int      { return len(m) }
func (m Items) Swap(a, b int) { m[a], m[b] = m[b], m[a] }
func (m Items) Less(a, b int) bool {
	return m[a].Text() < m[b].Text()
}
