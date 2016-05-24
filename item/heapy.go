package main

import (
	"container/heap"
	"fmt"
)

func main() {
	items := map[string]int{
		"banana": 3,
		"apple":  2,
		"pear":   5,
	}

	pq := make(priorityQueue, len(items))
	i := 0
	for value, priority := range items {
		pq[i] = &item{
			priority: priority,
			value:    value,
		} // can't use pop as internally, it modifies index.
		i++
	}

	heap.Init(&pq) // we want to update index once and for all.

	// Insert a new item and then modify its priority.
	itm := &item{
		value:    "orange",
		priority: 1,
	}
	heap.Push(&pq, itm)
	pq.update(itm, itm.value, 5)

	for len(pq) > 0 {
		pitm := heap.Pop(&pq).(*item)
		fmt.Printf("%0.2d:%s ", pitm.priority, pitm.value)
	}
	fmt.Println()
}

// an item is something we manage in a priority queue
type item struct {
	index    int // index is needed by update
	priority int
	value    string
}

// priorityQueue implements heap.Interface and holds pointers to items.
type priorityQueue []*item

func (pq priorityQueue) Len() int {
	return len(pq)
}

// we want Pop to give us highest, not lowest priority, hence gt.
func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	item := x.(*item)
	n := len(*pq)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	itm := old[n-1]
	itm.index = -1 // for safety
	*pq = old[0 : n-1]
	return itm
}

// update modifies the priority and value of an item in the queue.
func (pq *priorityQueue) update(itm *item, value string, priority int) {
	itm.value = value
	itm.priority = priority
	heap.Fix(pq, itm.index)
}
