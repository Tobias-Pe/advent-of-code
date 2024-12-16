package main

// An stateHeap is a min-stateHeap of ints.
type stateHeap []state

func (h stateHeap) Len() int           { return len(h) }
func (h stateHeap) Less(i, j int) bool { return h[i].cost < h[j].cost }
func (h stateHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *stateHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(state))
}

func (h *stateHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
