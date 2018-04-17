package main

import (
	"fmt"
	"sync"
	"time"
)

type ResultManager struct {
	parseIn, parseOut  chan *Backlink
	msgSent, errorSent chan bool
	queue, errorQueue  *Q
	wg                 *sync.WaitGroup
	logger             *Logger

	maxDepth, workersCount int
}

func (rm *ResultManager) Init(parseIn, parseOut chan *Backlink, maxDepth, workersCount int, queue, errorQueue *Q) {
	rm.parseIn = parseIn
	rm.parseOut = parseOut
	rm.maxDepth = maxDepth
	rm.workersCount = workersCount
	rm.queue = queue
	rm.errorQueue = errorQueue

	rm.wg = &sync.WaitGroup{}
	rm.msgSent = make(chan bool, workersCount)
	rm.errorSent = make(chan bool, workersCount)
	rm.logger = new(Logger)
	rm.logger.Init()
}

func (rm *ResultManager) Start() (*sync.WaitGroup, error) {
	rm.wg.Add(1)
	go rWorker(rm)
	go wWorker(rm)
	go eWorker(rm)

	return rm.wg, nil
}

func (rm *ResultManager) Kill() {
	rm.errorQueue.Clear()
	//close(rm.msgSent)
	//close(rm.errorSent)
	//close(rm.parseIn)
	rm.wg.Done()
}

func handle(rm *ResultManager, msg *Backlink) {
	rm.logger.Log("result", msg)
	Counter++
	if msg.Depth < rm.maxDepth {
		for i := 0; i < len(msg.BLList); i++ {
			rm.queue.Push(msg.BLList[i])
		}
	}

	rm.msgSent <- true
}

func eWorker(rm *ResultManager) {
	for range rm.errorSent {
		msg := rm.errorQueue.Pop()

		if msg.ReTry == 0 && ErrorRetryCounter == 0 {
			for i := 0; i < rm.workersCount; i++ {
				rm.parseIn <- StopMessage
			}
			break
		}
		if msg.ReTry == 0 {
			continue
		}
		time.Sleep(ERROR_RETRY_TIME)

		if msg.ReTry != MAX_RETRY_COUNT {
			ErrorRetryCounter--
		}
		msg.ReTry--
		if msg.ReTry == 0 {
			rm.logger.Log("error", msg)
		}
		rm.parseIn <- msg
	}
	rm.wg.Done()
}

func wWorker(rm *ResultManager) {
	for range rm.msgSent {
		for i := 0; i < rm.workersCount; i++ {
			msg := rm.queue.Pop()
			if msg == nil && i == 0 && ErrorRetryCounter == 0 {
				msg = StopMessage
			}
			if msg == nil {
				continue
			}
			rm.parseIn <- msg
		}
	}
}

func rWorker(rm *ResultManager) {
	for msg := range rm.parseOut {
		fmt.Println("result in", msg.Url)
		if msg == nil {
			continue
		}
		if msg.Error != nil {
			if msg.ReTry == MAX_RETRY_COUNT {
				ErrorRetryCounter += MAX_RETRY_COUNT - 1
			}
			rm.errorQueue.Push(msg)
			rm.errorSent <- true
			continue
		}

		handle(rm, msg)
	}
}
