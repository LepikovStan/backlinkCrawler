package main

import (
	"container/heap"
	"fmt"
	"sync"
)

type ResultManager struct {
	wCount, wNum, maxDepth int
	stopped                bool

	in, out chan *Backlink
	wg      *sync.WaitGroup
	pQueue  *PriorityQueue
	logger  *Logger
	queue   *Queue
}

func (rm *ResultManager) Init(options *Options) error {
	rm.in = options.in
	rm.out = options.out
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
	return rm.wg, nil
}

func (rm *ResultManager) Stop() {
	close(rm.out)
	fmt.Println("result manager stopped")
}

func (rm *ResultManager) Shutdown() {
	for i := 0; i < rm.wCount; i++ {
		rm.wg.Done()
	}
}

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
		BL := heap.Pop(rm.pQueue).(*Backlink)
		fmt.Println("send next", rm.maxDepth, BL.Depth, rm.queue.Len())
		//if rm.queue.Len() > 0 && BL.Depth < rm.maxDepth {
		//
		//}
		if BL.Depth <= rm.maxDepth {
			rm.out <- BL
		}
	}
}

func rmWorker(rm *ResultManager, i int) {
	fmt.Println(fmt.Sprintf("result worker %d started", i))
	for msg := range rm.in {
		if msg.Error != nil {
			//sendNextToCrawl(rm)
			fmt.Println("msg.TryCount", msg.TryCount)
			if msg.TryCount > 0 {
				msg.TryCount--
				rm.queue.Push(msg)
			}
			rm.logger.Log("error", fmtErrorLog(msg)...)
			continue
		}

		if msg.Depth > rm.maxDepth {
			fmt.Println("msg.Depth > rm.maxDepth", msg.Url)
		}

		for i := 0; i < len(msg.BLList); i++ {
			heap.Push(rm.pQueue, msg.BLList[i])
		}

		sendNextToCrawl(rm)
		rm.logger.Log("result", fmtResultLog(msg)...)
		fmt.Println(fmt.Sprintf("result worker %d, url %s", i, msg.Url), len(rm.out))
	}
	rm.wg.Done()
}
