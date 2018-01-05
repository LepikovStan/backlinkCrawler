package main

import (
	"flag"
	"fmt"
	"github.com/LepikovStan/backlinkCrawler/crawler"
	"github.com/LepikovStan/backlinkCrawler/lib/queue"
	"github.com/LepikovStan/backlinkCrawler/lib/writer"
	"github.com/LepikovStan/backlinkCrawler/parser"
	"log"
	"os"
	"sync"
	"time"
)

var workersCount int
var maxDepth int
var counter int

func initFlags() {
	flag.IntVar(&workersCount, "workers", 1, "")
	flag.IntVar(&maxDepth, "depth", 0, "")
	flag.Parse()
}

func crawlWorker(in chan queue.Backlink, out chan queue.Backlink, wg *sync.WaitGroup) {
	crwlr := crawler.New()
	for msg := range in {
		if msg.Depth > maxDepth {
			out <- msg
			break
		}
		body, err := crwlr.Crawl(msg.Url)
		if err != nil {
			msg.Error = err
			out <- msg
			continue
		}

		msg.Body = body
		out <- msg
	}
	wg.Done()
}

func parseWorker(in chan queue.Backlink, out chan queue.Backlink, wg *sync.WaitGroup, Q *queue.Q) {
	prsr := parser.New()
	for msg := range in {
		if msg.Depth > maxDepth {
			break
		}
		if msg.Error != nil {
			out <- msg
			Q.Set(Q.PopBuffer(workersCount - len(Q.GetChan())))
			continue
		}
		urlList, err := prsr.Parse(msg.Body)
		if err != nil {
			msg.Error = err
			out <- msg
			Q.Set(Q.PopBuffer(workersCount - len(Q.GetChan())))
			continue
		}
		msg.BLList = queue.TransformUrlToBacklink(urlList, msg.Depth+1)
		Q.SetBuffer(msg.BLList)
		Q.Set(Q.PopBuffer(workersCount - len(Q.GetChan())))
		out <- msg
	}
	wg.Done()
}

func resultHandler(in chan queue.Backlink, wg *sync.WaitGroup) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	resultDir := fmt.Sprintf("%s/%s/%d", dir, "results", time.Now().Unix())
	err = CreateDirIfNotExist(resultDir)
	if err != nil {
		log.Fatal(err)
	}
	wr := writer.New(resultDir)
	defer wr.Destroy()

	for msg := range in {
		if msg.Depth > maxDepth {
			break
		}
		counter++
		if msg.Error != nil {
			result := fmt.Sprintf("%s:\n\terror:%s\n", msg.Url, msg.Error)
			fmt.Println(result)
			wr.WriteError(result)
			continue
		}
		result := fmt.Sprintf("\n%s:\n", msg.Url)
		for _, backlink := range msg.BLList {
			result = fmt.Sprintf("%s\t%s\n", result, backlink.Url)
		}
		fmt.Println(result)
		wr.WriteResult(result)
	}
	wg.Done()
}

func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	fmt.Println("Start...")
	start := time.Now()
	initFlags()

	counter = 0
	crawledCh := make(chan queue.Backlink, workersCount)
	crawlWG := sync.WaitGroup{}

	parsedCh := make(chan queue.Backlink, workersCount)
	parseWG := sync.WaitGroup{}
	resultWG := sync.WaitGroup{}

	Q := queue.New(workersCount)
	Q.SetBuffer(queue.TransformUrlToBacklink(crawler.GetStartList(""), 0))
	Q.Set(Q.PopBuffer(workersCount - len(Q.GetChan())))

	crawlWG.Add(workersCount)
	for i := 0; i < workersCount; i++ {
		go crawlWorker(Q.GetChan(), crawledCh, &crawlWG)
	}
	parseWG.Add(workersCount)
	for i := 0; i < workersCount; i++ {
		go parseWorker(crawledCh, parsedCh, &parseWG, Q)
	}
	resultWG.Add(1)
	go resultHandler(parsedCh, &resultWG)

	crawlWG.Wait()
	close(crawledCh)
	parseWG.Wait()
	close(parsedCh)
	resultWG.Wait()
	fmt.Println("total", counter)

	end := time.Now()
	fmt.Println("\n")
	fmt.Println(end.Sub(start))
}
