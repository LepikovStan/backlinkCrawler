package main

import (
	"container/heap"
	//"context"
	"fmt"
	"sync"
	//"time"
	"time"
)

type ResultManager struct {
	wCount, wNum, maxDepth int
	stopped                bool

	in, out, eIn, eOut chan *Backlink
	wg                 *sync.WaitGroup
	pQueue             *PriorityQueue
	logger             *Logger
	queue              *Queue
}

func (rm *ResultManager) Init(options *Options) error {
	rm.in = options.in
	rm.out = options.out
	rm.eIn = options.eIn
	rm.eOut = options.eOut
	rm.wCount = options.wCount
	rm.wg = &sync.WaitGroup{}
	rm.maxDepth = options.maxDepth
	rm.pQueue = options.pQueue
	rm.queue = options.queue
	rm.logger = options.lg

	return nil
}

func (rm *ResultManager) Start() (*sync.WaitGroup, error) {
	fmt.Println("result manager started")
	rm.wg.Add(rm.wCount)
	for i := 0; i < rm.wCount; i++ {
		go rmWorker(rm, i)
	}
	rm.wg.Add(1)
	go rmErrorWorker(rm)
	return rm.wg, nil
}

func (rm *ResultManager) Stop() {
	close(rm.out)
	fmt.Println("result manager stopped")
}

func (rm *ResultManager) Shutdown() {}

func fmtResultLog(msg *Backlink) []string {
	result := make([]string, len(msg.BLList)+1)
	result[0] = msg.Url
	for i := 0; i < len(msg.BLList); i++ {
		result[i+1] = fmt.Sprintf("    %s", msg.BLList[i].Url)
	}
	return result
}

func fmtErrorLog(msg *Backlink) []string {
	return []string{
		msg.Url,
		fmt.Sprintf("    %s", msg.Error),
	}
}

func sendNextToCrawl(rm *ResultManager) {
	freeSpaceToCrawl := rm.wCount - len(rm.out)
	for i := 0; i < freeSpaceToCrawl; i++ {
		if rm.pQueue.Len() == 0 {
			rm.out <- StopMessage
			continue
		}

		BL := heap.Pop(rm.pQueue).(*Backlink)
		if BL == nil {
			return
		}

		rm.out <- BL
	}
}

func handleResult(rm *ResultManager, msg *Backlink) bool {
	if msg == nil {
		return true
	}

	if msg.Shutdown == true {
		return false
	}

	if msg.Error != nil {
		fmt.Println("msg.Error ->", msg.Error)
		msg.TryCount--
		rm.queue.Push(msg)
		sendNextToCrawl(rm)
		sendErrToCrawl(rm)
		rm.logger.Log("error", fmtErrorLog(msg)...)
		return true
	}

	if msg.BLList[0].Depth <= rm.maxDepth {
		for i := 0; i < len(msg.BLList); i++ {
			heap.Push(rm.pQueue, msg.BLList[i])
		}
	}
	//
	sendNextToCrawl(rm)
	rm.logger.Log("result", fmtResultLog(msg)...)
	return true
}

func sendErrToCrawl(rm *ResultManager) {
	freeSpaceToCrawl := 1 - len(rm.eOut)
	for i := 0; i < freeSpaceToCrawl; i++ {
		BL := rm.queue.Pop()

		if BL == nil {
			return
		}

		rm.eOut <- BL
	}
}
func handleError(rm *ResultManager, msg *Backlink) bool {
	if msg == nil {
		return true
	}
	fmt.Println("handleError", msg.TryCount, rm.queue.Len() == 0)

	if msg.TryCount == 0 && rm.queue.Len() == 0 {
		return false
	}

	if msg.TryCount == 0 {
		return true
	}

	msg.TryCount--
	rm.logger.Log("error", fmtErrorLog(msg)...)
	rm.queue.Push(msg)

	sendNextToCrawl(rm)
	time.Sleep(time.Second * 5)
	sendErrToCrawl(rm)
	return true
}

func rmErrorWorker(rm *ResultManager) {
	for msg := range rm.eIn {
		if ok := handleError(rm, msg); !ok {
			break
		}
	}
	fmt.Println("Stop error worker")
	rm.wg.Done()
}

func rmWorker(rm *ResultManager, i int) {
	fmt.Println(fmt.Sprintf("result worker %d started", i))
	for msg := range rm.in {
		if ok := handleResult(rm, msg); !ok {
			break
		}
	}
	rm.wg.Done()
}
