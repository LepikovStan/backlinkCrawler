package main

import "sync"

const Q_SIZE = 10

type Queue struct {
	arr                        []interface{}
	head, tail, arrSize, count int
	mu                         *sync.RWMutex
}

func (q *Queue) Push(item interface{}) {
	q.mu.Lock()
	if q.tail == q.arrSize {
		q.arrSize = q.tail - q.head + Q_SIZE
		newArr := make([]interface{}, q.arrSize)
		copy(newArr, q.arr[q.head:q.tail])
		q.arr = newArr
		q.tail = q.tail - q.head
		q.head = 0
	}
	q.arr[q.tail] = item
	q.tail++
	q.count++
	q.mu.Unlock()
}

func (q *Queue) Pop() interface{} {
	q.mu.Lock()
	headItem := q.arr[q.head]
	q.head++
	q.count--
	q.mu.Unlock()
	return headItem
}

func (q *Queue) Len() int {
	return q.count
}

func NewOueue() *Queue {
	q := new(Queue)
	q.arr = make([]interface{}, Q_SIZE)
	q.arrSize = Q_SIZE
	q.mu = &sync.RWMutex{}
	return q
}
