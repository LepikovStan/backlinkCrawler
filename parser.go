package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"sync"
)

// Parser is the type, contains the basic Parse method
type Parser struct {
	wCount, maxDepth, wNum int

	in, out chan *Backlink
	wg      *sync.WaitGroup
	queue   *Q
}

func (p *Parser) Init(options *Options) error {
	p.in = options.in
	p.out = options.out
	p.wCount = options.wCount
	p.maxDepth = options.maxDepth
	p.queue = options.queue

	p.wg = &sync.WaitGroup{}
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
	close(p.out)
	for i := 0; i < p.wCount; i++ {
		p.wg.Done()
	}
}

func pWorker(p *Parser, i int) {
	fmt.Println(fmt.Sprintf("parse worker %d started", i))
	for msg := range p.in {
		if msg.Error != nil {
			fmt.Println(fmt.Sprintf("parser in retry %s, attempt: %d", msg.Url, MAX_RETRY_COUNT-msg.ReTry))
		} else {
			fmt.Println("parser in", msg.Url, msg.Shutdown)
		}
		if msg == StopMessage {
			p.out <- msg
			break
		}

		BLList, err := parse(msg.Url)

		if err != nil {
			msg.Error = err
			p.out <- msg
			continue
		}

		msg.BLList = TransformUrlToBacklink(BLList, msg.Depth+1)
		p.out <- msg
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
