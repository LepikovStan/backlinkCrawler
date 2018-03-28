package main

import (
	"container/heap"
	"flag"
	"fmt"
	"github.com/LepikovStan/backlinkCrawler/lib"
	"log"
	"os"
	"os/signal"
	"sync"
)

var workersCount int
var maxDepth int

func parseFlags() {
	flag.IntVar(&workersCount, "workers", 2, "")
	flag.IntVar(&maxDepth, "depth", 0, "")
	flag.Parse()
}

type Worker interface {
	Start() (*sync.WaitGroup, error)
	Stop()
}

type Options struct {
	in, out          chan *Backlink
	wCount, maxDepth int
	lg               *Logger
	pQueue           *PriorityQueue
	queue            *Queue
}

type Backlink struct {
	Url             string
	Body            []byte
	BLList          []*Backlink
	Error           error
	Depth, TryCount int
	Shutdown        bool

	index, priority int
}

func getStartList(path string) []string {
	if path == "" {
		path = "./input.txt"
	}
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return lib.ReadFile(fmt.Sprintf("%s/%s", dir, path))
}

func TransformUrlToBacklink(urls []string, depth int) []*Backlink {
	result := make([]*Backlink, len(urls))
	for i, url := range urls {
		result[i] = &Backlink{
			Url:      url,
			Body:     nil,
			BLList:   nil,
			Error:    nil,
			Depth:    depth,
			TryCount: 3,
			priority: depth + 10,
		}
	}
	return result
}

func initializeCrawlIn(crawlIn chan *Backlink, pQueue *PriorityQueue) {
	startList := TransformUrlToBacklink(getStartList(""), 0)

	for i := 0; i < len(startList); i++ {
		heap.Push(pQueue, startList[i])
		//crawlIn <- startList[i]
	}
	for i := 0; i < workersCount; i++ {
		crawlIn <- heap.Pop(pQueue).(*Backlink)
	}
}

func interceptSignal(chIn chan *Backlink, ehwg *sync.WaitGroup) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	<-c
	fmt.Println()
	fmt.Println("Graceful shutdown...")
	ehwg.Done()
	for i := 0; i < workersCount; i++ {
		chIn <- &Backlink{
			Url:      "",
			Body:     nil,
			BLList:   nil,
			Error:    nil,
			Depth:    0,
			Shutdown: true,
			priority: 0,
		}
	}
}

func main() {
	parseFlags()
	chIn := make(chan *Backlink, workersCount)
	chOut := make(chan *Backlink, workersCount)
	pQueue := NewPQueue()
	queue := NewOueue()

	initializeCrawlIn(chIn, pQueue)
	lg := new(Logger)
	lg.Init()

	parser := new(Parser)
	parser.Init(&Options{
		chIn,
		chOut,
		workersCount,
		maxDepth,
		nil,
		nil,
		nil,
	})

	resultManager := new(ResultManager)
	resultManager.Init(&Options{
		chOut,
		chIn,
		workersCount,
		maxDepth,
		lg,
		pQueue,
		queue,
	})

	errorHandler := new(ErrorHandler)
	errorHandler.Init(&Options{
		nil,
		chIn,
		1,
		maxDepth,
		nil,
		pQueue,
		queue,
	})

	pwg, _ := parser.Start()
	rmwg, _ := resultManager.Start()
	ehwg, _ := errorHandler.Start()
	go interceptSignal(chIn, ehwg)

	fmt.Println()

	ehwg.Wait()
	errorHandler.Stop()

	pwg.Wait()
	parser.Stop()

	rmwg.Wait()
	resultManager.Stop()
}
