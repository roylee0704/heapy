package main

import (
	"container/heap"
	"fmt"
)

// this is a sandbox environment for heap
//
// 1. make a pool.
// 2. attached pool with heap apis.
// 3. pop and push, let hash do the hard work

// An IntHeap is a min-heap of ints.
type intHeap []int

func (h intHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h intHeap) Len() int {
	return len(h)
}

func (h intHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *intHeap) Push(val interface{}) {
	*h = append(*h, val.(int))
}

func (h *intHeap) Pop() interface{} {

	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func main() {

	poll := &intHeap{2, 1, 4, 5, 3}
	fmt.Println("before: ", poll)

	heap.Init(poll)
	fmt.Println("after: ", poll)

	for len(*poll) > 0 {
		fmt.Println(heap.Pop(poll))
	}
}
