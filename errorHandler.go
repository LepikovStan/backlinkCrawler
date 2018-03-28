package main

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

type ErrorHandler struct {
	wCount, maxDepth, wNum int
	paused                 bool

	out    chan *Backlink
	wg     *sync.WaitGroup
	queue  *Queue
	pQueue *PriorityQueue
}

func (eh *ErrorHandler) Init(options *Options) {
	eh.out = options.out
	eh.wCount = options.wCount
	eh.wg = &sync.WaitGroup{}
	eh.maxDepth = options.maxDepth
	eh.queue = options.queue
	eh.pQueue = options.pQueue
}
func (eh *ErrorHandler) Start() (*sync.WaitGroup, error) {
	fmt.Println("error handler started")
	eh.wg.Add(eh.wCount)
	for i := 0; i < eh.wCount; i++ {
		eh.wNum = i
		go ehWorker(eh, i)
	}
	return eh.wg, nil
}
func (eh *ErrorHandler) Stop() {
	fmt.Println("error handler stopped")
}

func (eh *ErrorHandler) Pause() {
	eh.paused = true
}

func (eh *ErrorHandler) Continue() {
	eh.paused = false
	timeout(eh)
}

func timeout(eh *ErrorHandler) {
	time.Sleep(time.Second * 5)
	if eh.paused == true {
		return
	}
	msg := eh.queue.Pop().(*Backlink)
	fmt.Println("timeout", msg)
	if msg == nil {
		timeout(eh)
		return
	}
	if len(eh.out) > 0 {
		msg.priority = 0
		heap.Push(eh.pQueue, msg)
		return
	}
	eh.out <- msg
	timeout(eh)
}

func ehWorker(eh *ErrorHandler, wNum int) {
	timeout(eh)
}
