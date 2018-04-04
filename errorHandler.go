package main

//
//import (
//	"container/heap"
//	"fmt"
//	"sync"
//	"time"
//)
//
//type ErrorHandler struct {
//	wCount, maxDepth, wNum int
//
//	out    chan *Backlink
//	wg     *sync.WaitGroup
//	queue  *Queue
//	pQueue *PriorityQueue
//}
//
//func (eh *ErrorHandler) Init(options *Options) {
//	eh.out = options.out
//	eh.wCount = options.wCount
//	eh.wg = &sync.WaitGroup{}
//	eh.maxDepth = options.maxDepth
//	eh.queue = options.queue
//	eh.pQueue = options.pQueue
//}
//func (eh *ErrorHandler) Start() (*sync.WaitGroup, error) {
//	fmt.Println("error handler started")
//	eh.wg.Add(eh.wCount)
//	for i := 0; i < eh.wCount; i++ {
//		eh.wNum = i
//		go ehWorker(eh, i)
//	}
//	return eh.wg, nil
//}
//func (eh *ErrorHandler) Stop() {
//	fmt.Println("error handler stopped")
//}
//
//func (eh *ErrorHandler) Shutdown() {
//	for i := 0; i < eh.wCount; i++ {
//		eh.wg.Done()
//	}
//}
//
//func ehWorker(eh *ErrorHandler, wNum int) {
//	for msg := range eh.queue.GetChan() {
//		fmt.Println("msg.Shutdown", msg.Shutdown)
//		if msg.Shutdown == true {
//			break
//		}
//
//		time.Sleep(time.Second * 5)
//		if len(eh.out) > 0 {
//			msg.priority = 0
//			heap.Push(eh.pQueue, msg)
//			continue
//		}
//
//		eh.out <- msg
//	}
//	eh.wg.Done()
//	eh.wCount--
//	fmt.Println("error handler shutdown", wNum)
//}
