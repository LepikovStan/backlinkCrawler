package main

import (
	"fmt"
	"sync"
)

const Q_SIZE = 10

type Queue struct {
	arr                        []*Backlink
	head, tail, arrSize, count int
	mu                         *sync.RWMutex
}

func (q *Queue) Push(item *Backlink) {
	q.mu.Lock()
	if q.tail == q.arrSize {
		q.arrSize = q.tail - q.head + Q_SIZE
		newArr := make([]*Backlink, q.arrSize)
		copy(newArr, q.arr[q.head:q.tail])
		q.arr = newArr
		q.tail = q.tail - q.head
		q.head = 0
	}
	q.arr[q.tail] = item
	q.tail++
	q.count++
	fmt.Println("queue push", q.tail, q.head, q.count)
	q.mu.Unlock()
}

func (q *Queue) Pop() *Backlink {
	q.mu.Lock()
	headItem := q.arr[q.head]
	q.head++
	q.count--
	if q.count < 0 {
		q.count = 0
		q.tail = q.head
	}
	fmt.Println("queue pop", q.tail, q.head, q.count)
	q.mu.Unlock()
	return headItem
}

func (q *Queue) Len() int {
	fmt.Println("queue len", q.tail, q.head, q.count)
	return q.tail - q.head
}

func NewOueue() *Queue {
	q := new(Queue)
	q.arr = make([]*Backlink, Q_SIZE)
	q.arrSize = Q_SIZE
	q.mu = &sync.RWMutex{}
	return q
}
