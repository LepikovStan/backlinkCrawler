package main

import (
	//"container/heap"
	"flag"
	//"fmt"
	//"github.com/LepikovStan/backlinkCrawler/lib"
	//"log"
	//"os"
	//"os/signal"
	"fmt"
	"os"
	"os/signal"
	"sync"
)

var workersCount int

//var errorWorkersCount int
var maxDepth int

var StopMessage = &Backlink{
	Shutdown: true,
}

func parseFlags() {
	flag.IntVar(&workersCount, "workers", 2, "")
	//flag.IntVar(&errorWorkersCount, "errorWorkers", 1, "")
	flag.IntVar(&maxDepth, "depth", 0, "")
	flag.Parse()
}

type Worker interface {
	Start() (*sync.WaitGroup, error)
	Close()
	Shutdown()
	Kill()
}

type Options struct {
	wCount, ewCount, maxDepth int

	in, out, eIn, eOut chan *Backlink
	lg                 *Logger
	pQueue, epQueue    *PriorityQueue
	queue              *Q
}

type Backlink struct {
	Url          string
	Body         []byte
	BLList       []*Backlink
	Error        error
	Depth, ReTry int
	Shutdown     bool

	index, priority int
}

//func getStartList(path string) []string {
//	if path == "" {
//		path = "./input.txt"
//	}
//	dir, err := os.Getwd()
//	if err != nil {
//		log.Fatal(err)
//	}
//	return lib.ReadFile(fmt.Sprintf("%s/%s", dir, path))
//}

func TransformUrlToBacklink(urls []string, depth int) []*Backlink {
	result := make([]*Backlink, len(urls))
	for i, url := range urls {
		result[i] = &Backlink{
			Url:    url,
			Body:   nil,
			BLList: nil,
			Error:  nil,
			Depth:  depth,
			ReTry:  2,
			//priority: depth + 10,
		}
	}
	return result
}

func interceptSignal(w ...Worker) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	<-c
	fmt.Println()
	fmt.Println("Force Shutdown...")
	for i := 0; i < len(w); i++ {
		w[i].Shutdown()
	}
}

var Counter int

func main() {
	parseFlags()
	parseIn := make(chan *Backlink, workersCount)
	parseOut := make(chan *Backlink, workersCount)
	errorParseIn := make(chan *Backlink, workersCount)
	errorParseOut := make(chan *Backlink, workersCount)
	queue := NewQ()
	errorQueue := NewQ()

	startList := TransformUrlToBacklink([]string{"golang.org"}, 0)
	for i := 0; i < len(startList); i++ {
		parseIn <- startList[i]
	}

	parser := new(Parser)
	parser.Init(&Options{
		in:       parseIn,
		out:      parseOut,
		wCount:   workersCount,
		maxDepth: maxDepth,
		queue:    queue,
		//pQueue:   pQueue,
	})

	errorsParser := new(Parser)
	errorsParser.Init(&Options{
		in:       errorParseIn,
		out:      errorParseOut,
		wCount:   workersCount,
		maxDepth: maxDepth,
		queue:    errorQueue,
		//pQueue:   pQueue,
	})

	resultManager := new(ResultManager)
	resultManager.Init(parseIn, parseOut, errorParseIn, errorParseOut, maxDepth, workersCount, queue, errorQueue)

	pwg, _ := parser.Start()
	errorsParser.Start()
	resultManager.Start()
	//go interceptSignal(parser)

	pwg.Wait()
	parser.Close()
	fmt.Println("Total ->", Counter)
}

//func initializeCrawlIn(crawlIn chan *Backlink, pQueue *PriorityQueue) {
//	startList := TransformUrlToBacklink(getStartList(""), 0)
//	msgToStart := len(startList)
//	if workersCount <= msgToStart {
//		msgToStart = workersCount
//	}
//
//	for i := 0; i < len(startList); i++ {
//		heap.Push(pQueue, startList[i])
//		//crawlIn <- startList[i]
//	}
//	fmt.Println("msgToStart", msgToStart)
//	for i := 0; i < msgToStart; i++ {
//		crawlIn <- heap.Pop(pQueue).(*Backlink)
//	}
//}

//func main() {
//	parseFlags()
//	parseIn := make(chan *Backlink, workersCount)
//	parseOut := make(chan *Backlink, workersCount)
//	errorsIn := make(chan *Backlink, errorWorkersCount)
//	errorsOut := make(chan *Backlink, errorWorkersCount)
//	pQueue := NewPQueue()
//	epQueue := NewPQueue()
//	//queue := NewOueue()
//	lg := new(Logger)
//	lg.Init()
//
//	initializeCrawlIn(parseIn, pQueue)
//
//	parser := new(Parser)
//	parser.Init(&Options{
//		in:       parseIn,
//		out:      parseOut,
//		wCount:   workersCount,
//		maxDepth: maxDepth,
//		pQueue:   pQueue,
//	})
//
//	eParser := new(Parser)
//	eParser.Init(&Options{
//		in:       errorsIn,
//		out:      errorsOut,
//		wCount:   workersCount,
//		maxDepth: maxDepth,
//		pQueue:   epQueue,
//	})
//
//	resultManager := new(ResultManager)
//	resultManager.Init(&Options{
//		in:       parseOut,
//		out:      parseIn,
//		eIn:      errorsOut,
//		eOut:     errorsIn,
//		wCount:   workersCount,
//		ewCount:  errorWorkersCount,
//		maxDepth: maxDepth,
//		lg:       lg,
//		pQueue:   pQueue,
//		epQueue:  epQueue,
//	})
//
//	pwg, _ := parser.Start()
//	rmwg, _ := resultManager.Start()
//	eParser.Start()
//	go interceptSignal(eParser, parser, resultManager)
//
//	fmt.Println()
//
//	fmt.Println("done 2")
//	pwg.Wait()
//	parser.Stop()
//
//	fmt.Println("done 3")
//	rmwg.Wait()
//	resultManager.Stop()
//	fmt.Println("done 4")
//}
