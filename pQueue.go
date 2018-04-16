// This example demonstrates a priority queue built using the heap interface.
package main

import (
	"container/heap"
	//"fmt"
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
	//fmt.Println("rm.epQueue Pust")
	n := len(*pq)
	item := x.(*Backlink)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	//fmt.Println("rm.epQueue Pop", old)
	n := len(old)
	//fmt.Println("rm.epQueue Pop", n)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	//fmt.Println("rm.epQueue Pop", *pq)
	return item
}

func NewPQueue() *PriorityQueue {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	return &pq
}
