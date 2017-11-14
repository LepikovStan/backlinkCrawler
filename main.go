package main

import (
	"flag"
	"fmt"
	"github.com/LepikovStan/backlinkCrawler/crawler"
	"github.com/LepikovStan/backlinkCrawler/lib/queue"
	"github.com/LepikovStan/backlinkCrawler/parser"
	//"./crawler"
	//"./lib/queue"
	//"./parser"
	"time"
	// "sync"
	"sync"
)

var parsersCount int
var crawlersCount int
var maxDepth int

func getStartList() []string {
	return []string{
		// "https://yandex.ru/yandsearch?text=golang%20error%20type&lr=2",
		// "https://gobyexample.com/errors",
		// "https://golang.org/pkg/sync",
		//"https://developers.google.com/products/",
		"http://www.tattyworld.net",
	}
}

func initFlags() {
	flag.IntVar(&parsersCount, "parsers", 1, "")
	flag.IntVar(&crawlersCount, "crawlers", 1, "")
	flag.IntVar(&maxDepth, "depth", 1, "")
	flag.Parse()
}

func crawlWorker(in chan queue.Backlink, out chan queue.Backlink, wg *sync.WaitGroup) {
	crwlr := crawler.New()
	for msg := range in {
		if msg.Depth > maxDepth {
			break
		}
		body, err := crwlr.Crawl(msg.Url)
		//fmt.Println("crawl depth", msg.Depth, maxDepth)
		if err != nil {
			msg.Error = err
			fmt.Println("crawl worker error", err)
			out <- msg
			continue
		}

		msg.Body = body
		out <- msg
	}
	//fmt.Println("crawlWorker Done -<<<<<<<<<<<<<<<<<")
	wg.Done()
}

func parseWorker(in chan queue.Backlink, out chan queue.Backlink, wg *sync.WaitGroup, Q *queue.Q) {
	prsr := parser.New()
	for msg := range in {
		if msg.Error != nil {
			out <- msg
			emptyBacklink := queue.Backlink{
				Depth: msg.Depth+1,
			}
			Q.SetBuffer([]queue.Backlink{emptyBacklink})
			Q.Write(Q.PopBuffer(parsersCount - len(in)))
			continue
		}
		urlList, err := prsr.Parse(msg.Body)
		if err != nil {
			msg.Error = err
			//fmt.Println("parse worker error", err)
			out <- msg
			continue
		}
		msg.BLList = transformUrlToBacklink(urlList, msg.Depth+1)
		Q.SetBuffer(msg.BLList)
		Q.Write(Q.PopBuffer(parsersCount - len(in)))
		out <- msg
	}
	//fmt.Println("parseWorker Done")
	wg.Done()
}

func resultHandler(in chan queue.Backlink) {
	for msg := range in {
		if msg.Error != nil {
			fmt.Println("resultHandler error ->", msg.Error)
		}
		fmt.Println("resultHandler msg", msg.Url, msg.Depth)
	}
}

func transformUrlToBacklink(urls []string, depth int) []queue.Backlink {
	result := make([]queue.Backlink, len(urls))
	for i, url := range urls {
		result[i] = queue.Backlink{
			Url:    url,
			Body:   nil,
			BLList: nil,
			Error:  nil,
			Depth:  depth,
		}
	}
	return result
}

func main() {
	fmt.Println("Start...")
	start := time.Now()
	initFlags()

	crawledCh := make(chan queue.Backlink, parsersCount)
	crawlWG := sync.WaitGroup{}

	parsedCh := make(chan queue.Backlink, parsersCount)
	parseWG:= sync.WaitGroup{}

	Q := queue.New(parsersCount)
	Q.Write(transformUrlToBacklink(getStartList(), 0))

	crawlWG.Add(crawlersCount)
	for i := 0; i < crawlersCount; i++ {
		go crawlWorker(Q.GetChan(), crawledCh, &crawlWG)
	}
	parseWG.Add(parsersCount)
	for i := 0; i < parsersCount; i++ {
		go parseWorker(crawledCh, parsedCh, &parseWG, Q)
	}
	go resultHandler(parsedCh)

	crawlWG.Wait()
	close(crawledCh)
	parseWG.Wait()

	end := time.Now()
	fmt.Println("\n")
	fmt.Println(end.Sub(start))
}
