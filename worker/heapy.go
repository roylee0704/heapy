package main

import (
	"bytes"
	"container/heap"
	"flag"
	"fmt"
	"time"
)

var (
	nWorker  int
	nRequest int
)

func init() {
	flag.IntVar(&nWorker, "w", 5, "number of workers")
	flag.IntVar(&nRequest, "r", 100, "number of requests")
}

func main() {

	flag.Parse()
	b := NewBalancer(nWorker)

	b.simDispatch()
	b.simCompletion()

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

// Complete allows worker to report for job completion.
func (b Balancer) Complete(w *Worker) {
	heap.Remove(&b.pool, w.i)
	w.pending--
	heap.Push(&b.pool, w)
}

// Dispatch distributes job to most ligthly loaded worker.
func (b Balancer) Dispatch() {
	w := heap.Pop(&b.pool).(*Worker)
	w.pending++
	heap.Push(&b.pool, w)
}

// Print logs # of pending tasks for each workers.
func (b Balancer) Print() {
	var buffer bytes.Buffer

	var sum int
	for _, worker := range b.pool {
		buffer.WriteString(fmt.Sprintf("%3d", worker.pending))
		sum += worker.pending
	}
	avg := float64(sum) / float64(nWorker)
	fmt.Printf("%s   %.2f\n", buffer.String(), avg)
}

// simDispatch naively simulates job requests for each worker
func (b Balancer) simDispatch() {
	for i := 0; i < nRequest; i++ {
		b.Dispatch()
		b.Print()
	}
}

// simCompletion naively simulates job completion(all) for each worker
func (b Balancer) simCompletion() {

	for empty := nWorker; empty > 0; {
		for _, w := range b.pool {
			time.Sleep(100 * time.Millisecond) // 1 second
			if w.pending > 0 {
				b.Complete(w)
				b.Print()
				continue
			}
			empty--
		}
	}

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
