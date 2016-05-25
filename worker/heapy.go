package main

import (
	"bytes"
	"container/heap"
	"flag"
	"fmt"
)

var (
	nWorkers = flag.Int("n", 5, "number of workers")
)

func main() {

	flag.Parse()
	b := NewBalancer(*nWorkers)

	for i := 0; i < 10; i++ {
		w := heap.Pop(&b.pool).(*Worker)
		w.pending++
		heap.Push(&b.pool, w)

		b.Print()
	}

}

// Balancer dispatches jobs to workers, and update upon completion.
type Balancer struct {
	pool Pool
}

// NewBalancer returns new balancers with a pool of n workers.
func NewBalancer(n int) *Balancer {
	return &Balancer{genWorker(n)}
}

// genWorker produces workers based on number to generate.
func genWorker(n int) Pool {
	p := make(Pool, 0, n)
	for i := 0; i < n; i++ {
		w := &Worker{
			pending: 0,
			name:    fmt.Sprintf("[%.2d]", i),
		}
		heap.Push(&p, w)
	}
	return p
}

// Print logs # of pending tasks for each workers.
func (b Balancer) Print() {
	var buffer bytes.Buffer
	for _, worker := range b.pool {
		buffer.WriteString(fmt.Sprintf("%2d", worker.pending))
	}
	fmt.Println(buffer.String())
}

// Worker is something we manage in a pending queue
type Worker struct {
	i       int // index is needed by update
	pending int
	name    string
}

// Pool implements heap.Interface and holds pointers to Workers.
type Pool []*Worker

func (p Pool) Len() int { return len(p) }

func (p Pool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

func (p *Pool) Swap(i, j int) {
	a := *p
	a[i], a[j] = a[j], a[i]
	a[i].i = i
	a[j].i = j
}

func (p *Pool) Push(x interface{}) {
	a := *p
	n := len(a)
	a = a[0 : n+1]
	w := x.(*Worker)
	a[n] = w
	w.i = n
	*p = a
}

func (p *Pool) Pop() interface{} {
	a := *p
	*p = a[0 : len(a)-1]
	w := a[len(a)-1]
	w.i = -1 // for safety
	return w
}
