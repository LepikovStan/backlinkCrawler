package main

import (
	//"container/heap"
	//"context"
	//"fmt"
	"sync"
	//"time"
	//"time"
	//"fmt"
	"fmt"
)

type ResultManager struct {
	parseIn, parseOut, errorParseIn, errorParseOut chan *Backlink
	queue, errorQueue                              *Q
	maxDepth, workersCount                         int
	//wCount, wNum, maxDepth int
	//stopped                bool
	//
	//in, out, eIn, eOut chan *Backlink
	wg *sync.WaitGroup
	//pQueue, epQueue    *PriorityQueue
	//logger             *Logger
}

func (rm *ResultManager) Init(parseIn, parseOut, errorParseIn, errorParseOut chan *Backlink, maxDepth, workersCount int, queue, errorQueue *Q) {
	rm.parseIn = parseIn
	rm.parseOut = parseOut
	rm.errorParseIn = errorParseIn
	rm.errorParseOut = errorParseOut

	rm.maxDepth = maxDepth
	rm.workersCount = workersCount
	rm.queue = queue
	rm.errorQueue = errorQueue
	//rm.in = options.in
	//rm.out = options.out
	//rm.eIn = options.eIn
	//rm.eOut = options.eOut
	//rm.wCount = options.wCount
	rm.wg = &sync.WaitGroup{}
	//rm.maxDepth = options.maxDepth
	//rm.pQueue = options.pQueue
	//rm.epQueue = options.epQueue
	//rm.logger = options.lg
}

func (rm *ResultManager) Start() (*sync.WaitGroup, error) {
	rm.wg.Add(1)
	go worker(rm)
	//fmt.Println("result manager started")
	//rm.wg.Add(rm.workersCount)
	//for i := 0; i < rm.workersCount; i++ {
	//	go worker(rm)
	//}
	//rm.wg.Add(1)
	//go rmErrorWorker(rm)
	return rm.wg, nil
}

func handle(rm *ResultManager, msg *Backlink) {
	if msg.Depth < rm.maxDepth {
		for i := 0; i < len(msg.BLList); i++ {
			rm.queue.Push(msg.BLList[i])
		}
	}

	parserFreeSpace := rm.workersCount - len(rm.parseIn)
	for i := 0; i < parserFreeSpace; i++ {
		m := rm.queue.Pop()
		if m == nil && i == 0 {
			m = StopMessage
		}
		if m == nil && i != 0 {
			continue
		}
		rm.parseIn <- m
	}
}

func handleError(rm *ResultManager, msg *Backlink) {
	fmt.Println(msg.ReTry)
	rm.errorQueue.Push(msg)

	parserFreeSpace := rm.workersCount - len(rm.errorParseIn)
	fmt.Println(parserFreeSpace)
	for i := 0; i < parserFreeSpace; i++ {
		m := rm.errorQueue.Pop()
		if m == nil && i == 0 {
			m = StopMessage
		}
		if m == nil && i != 0 {
			continue
		}
		fmt.Println("ae")
		rm.errorParseIn <- m
	}
}

func worker(rm *ResultManager) {
	for msg := range rm.parseOut {
		fmt.Println("result in", msg.Url)
		if msg == nil {
			continue
		}
		if msg.Error != nil {
			handleError(rm, msg)
			continue
		}

		handle(rm, msg)

		//if ok := handleResult(rm, msg); !ok {
		//	break
		//}
	}
	rm.wg.Done()
}

//func (rm *ResultManager) Stop() {
//	close(rm.out)
//	close(rm.eOut)
//	fmt.Println("result manager stopped")
//}
//
//func (rm *ResultManager) Shutdown() {}
//
//func fmtResultLog(msg *Backlink) []string {
//	result := make([]string, len(msg.BLList)+1)
//	result[0] = msg.Url
//	for i := 0; i < len(msg.BLList); i++ {
//		result[i+1] = fmt.Sprintf("    %s", msg.BLList[i].Url)
//	}
//	return result
//}
//
//func fmtErrorLog(msg *Backlink) []string {
//	return []string{
//		msg.Url,
//		fmt.Sprintf("    %s", msg.Error),
//	}
//}
//
//func sendNextToCrawl(rm *ResultManager) {
//	freeSpaceToCrawl := rm.wCount - len(rm.out)
//	for i := 0; i < freeSpaceToCrawl; i++ {
//		if rm.pQueue.Len() == 0 {
//			rm.out <- StopMessage
//			continue
//		}
//
//		BL := heap.Pop(rm.pQueue).(*Backlink)
//		if BL == nil {
//			return
//		}
//
//		rm.out <- BL
//	}
//}
//
//func handleResult(rm *ResultManager, msg *Backlink) bool {
//	if msg == nil {
//		return true
//	}
//
//	if msg.Shutdown == true {
//		return false
//	}
//
//	if msg.Error != nil {
//		fmt.Println("msg.Error ->", msg.Error)
//		msg.TryCount--
//		fmt.Println("result rm.epQueue push", msg.TryCount)
//		heap.Push(rm.epQueue, msg)
//		sendNextToCrawl(rm)
//		rm.eIn <- msg
//		//sendErrToCrawl(rm)
//		rm.logger.Log("error", fmtErrorLog(msg)...)
//		return true
//	}
//
//	if msg.BLList[0].Depth <= rm.maxDepth {
//		for i := 0; i < len(msg.BLList); i++ {
//			heap.Push(rm.pQueue, msg.BLList[i])
//		}
//	}
//	//
//	sendNextToCrawl(rm)
//	rm.logger.Log("result", fmtResultLog(msg)...)
//	return true
//}
//
//func sendErrToCrawl(rm *ResultManager) {
//	freeSpaceToCrawl := 1 - len(rm.eOut)
//	for i := 0; i < freeSpaceToCrawl; i++ {
//		BL := heap.Pop(rm.epQueue).(*Backlink)
//		fmt.Println("rm.epQueue pop", BL.TryCount)
//
//		if BL == nil {
//			return
//		}
//
//		rm.eOut <- BL
//	}
//}
//func handleError(rm *ResultManager, msg *Backlink) bool {
//	if msg == nil {
//		return true
//	}
//	msg.TryCount--
//	fmt.Println("handleError rm.epQueue", msg.TryCount, rm.epQueue.Len() == 0, rm.epQueue)
//
//	if msg.TryCount == 0 && rm.epQueue.Len() == 0 {
//		return false
//	}
//
//	if msg.TryCount == 0 {
//		return true
//	}
//
//	rm.logger.Log("error", fmtErrorLog(msg)...)
//	fmt.Println("rm.epQueue push", msg.TryCount)
//	heap.Push(rm.epQueue, msg)
//
//	sendNextToCrawl(rm)
//	time.Sleep(time.Second * 5)
//	sendErrToCrawl(rm)
//	return true
//}
//
//func rmErrorWorker(rm *ResultManager) {
//	for msg := range rm.eIn {
//		if ok := handleError(rm, msg); !ok {
//			break
//		}
//	}
//	fmt.Println("Stop error worker")
//	rm.wg.Done()
//}
//
//func rmWorker(rm *ResultManager, i int) {
//	fmt.Println(fmt.Sprintf("result worker %d started", i))
//	for msg := range rm.in {
//		if ok := handleResult(rm, msg); !ok {
//			break
//		}
//	}
//	rm.wg.Done()
//}
