package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"sync"
	//"time"
	//"container/heap"
)

// Parser is the type, contains the basic Parse method
type Parser struct {
	wCount, maxDepth, wNum int

	in, out chan *Backlink
	wg      *sync.WaitGroup
	pQueue  *PriorityQueue
	queue   *Q
}

func (p *Parser) Init(options *Options) error {
	p.in = options.in
	p.out = options.out
	p.wCount = options.wCount
	p.wg = &sync.WaitGroup{}
	p.maxDepth = options.maxDepth
	p.pQueue = options.pQueue
	p.queue = options.queue
	return nil
}

func (p *Parser) Start() (*sync.WaitGroup, error) {
	//fmt.Println("parser started")
	p.wg.Add(p.wCount)
	for i := 0; i < p.wCount; i++ {
		p.wNum = i
		go pWorker(p, i)
	}
	return p.wg, nil
}

func (p *Parser) Close() {
	close(p.out)
	//fmt.Println("parser stopped")
}

func (p *Parser) Shutdown() {
	for i := 0; i < p.wCount; i++ {
		p.queue.Unshift(StopMessage)
	}
}

func (p *Parser) Kill() {
	for i := 0; i < p.wCount; i++ {
		p.wg.Done()
	}
}

func pWorker(p *Parser, i int) {
	fmt.Println(fmt.Sprintf("parse worker %d started", i))
	for msg := range p.in {
		fmt.Println("parser in", msg.Url, msg.Shutdown)
		//startTime := time.Now()
		if msg.Shutdown == true {
			p.out <- msg
			break
		}
		//if msg.Depth > p.maxDepth {
		//	break
		//}
		//fmt.Println(fmt.Sprintf("parse worker %d parse url %s, depth %d", i, msg.Url, msg.Depth))
		if msg.Error != nil {
			p.out <- msg
			continue
		}
		BLList, err := parse(msg.Url)

		if err != nil {
			msg.Error = err
			p.out <- msg
			continue
		}

		msg.BLList = TransformUrlToBacklink(BLList, msg.Depth+1)
		Counter++
		p.out <- msg
		//fmt.Println(fmt.Sprintf("parse worker %d parsed url %s, %s", i, msg.Url, time.Now().Sub(startTime)))
		//fmt.Println("result", urlList)
	}
	p.wg.Done()
	p.wCount--
	fmt.Println("parser shutdown", i)
}

// Parse parsing html and find all links on the page
// return list of finded links or error
func parse(url string) ([]string, error) {
	doc, err := goquery.NewDocument(url)
	result := []string{}

	if err != nil {
		fmt.Println("Parser error", err)
		return nil, err
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		Href, _ := s.Attr("href")
		if !strings.HasPrefix(Href, "http") {
			return
		}
		result = append(result, Href)
	})

	return result, nil
}
