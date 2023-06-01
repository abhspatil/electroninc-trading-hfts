package utils

import dtov1 "github.com/abhspatil/electronic-trading/pkg/dto/v1"

type OrderHeap []*dtov1.Order

func (h OrderHeap) Len() int { return len(h) }

func (h OrderHeap) Less(i, j int) bool {
	if h[i].Price < h[j].Price {
		return true
	} else if h[i].Price > h[j].Price {
		return false
	}
	return h[i].Quantity < h[j].Quantity
}

func (h OrderHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *OrderHeap) Push(x interface{}) {
	*h = append(*h, x.(*dtov1.Order))
}

func (h *OrderHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
