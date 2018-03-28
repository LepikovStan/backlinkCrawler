// This example demonstrates a priority queue built using the heap interface.
package main

import (
	"container/heap"
	//"fmt"
	"fmt"
)

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Backlink

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Backlink)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) Clear() {
	emptyPQ := make(PriorityQueue, 0)
	fmt.Println("clear", emptyPQ)
	pq = &emptyPQ
	fmt.Println("clear", pq)
}

func NewPQueue() *PriorityQueue {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	return &pq
}
