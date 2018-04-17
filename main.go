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
	"github.com/LepikovStan/backlinkCrawler/lib"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

const MAX_RETRY_COUNT = 2
const ERROR_RETRY_TIME = time.Second * 2

var workersCount int
var errorRetryCount int

//var errorWorkersCount int
var maxDepth int

var StopMessage = &Backlink{
	Shutdown: true,
}

func parseFlags() {
	flag.IntVar(&workersCount, "workers", 2, "")
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
			Url:    url,
			Body:   nil,
			BLList: nil,
			Error:  nil,
			Depth:  depth,
			ReTry:  MAX_RETRY_COUNT,
			//priority: depth + 10,
		}
	}
	return result
}

func interceptSignal(rm *ResultManager, w ...Worker) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	<-c
	fmt.Println()
	fmt.Println("Force Shutdown...")
	for i := 0; i < len(w); i++ {
		w[i].Shutdown()
	}
	rm.Kill()
}

var Counter int
var ErrorRetryCounter int

func main() {
	parseFlags()
	parseIn := make(chan *Backlink, workersCount)
	parseOut := make(chan *Backlink, workersCount)
	queue := NewQ()
	errorQueue := NewQ()

	startList := TransformUrlToBacklink(getStartList(""), 0)
	for i := 0; i < len(startList); i++ {
		queue.Push(startList[i])
	}
	for i := 0; i < workersCount; i++ {
		parseIn <- queue.Pop()
	}

	parser := new(Parser)
	parser.Init(&Options{
		in:       parseIn,
		out:      parseOut,
		wCount:   workersCount,
		maxDepth: maxDepth,
		queue:    queue,
	})

	resultManager := new(ResultManager)
	resultManager.Init(parseIn, parseOut, maxDepth, workersCount, queue, errorQueue)

	wg, _ := parser.Start()
	rwg, _ := resultManager.Start()
	go interceptSignal(resultManager, parser)

	wg.Wait()
	parser.Close()

	rwg.Wait()
	//fmt.Println("Total ->", Counter)
}
