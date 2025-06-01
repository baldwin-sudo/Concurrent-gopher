package main

// An IntHeap is a min-heap of ints.
type WorkersPool []*Worker

func (h WorkersPool) Len() int           { return len(h) }
func (h WorkersPool) Less(i, j int) bool { return h[i].Load < h[j].Load }
func (h WorkersPool) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *WorkersPool) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*Worker))
}

func (h *WorkersPool) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
