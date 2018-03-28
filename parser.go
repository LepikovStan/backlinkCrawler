package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"sync"
	//"time"
)

// Parser is the type, contains the basic Parse method
type Parser struct {
	wCount, maxDepth, wNum int

	in, out chan *Backlink
	wg      *sync.WaitGroup
}

func (p *Parser) Init(options *Options) error {
	p.in = options.in
	p.out = options.out
	p.wCount = options.wCount
	p.wg = &sync.WaitGroup{}
	p.maxDepth = options.maxDepth
	return nil
}

func (p *Parser) Start() (*sync.WaitGroup, error) {
	fmt.Println("parser started")
	p.wg.Add(p.wCount)
	for i := 0; i < p.wCount; i++ {
		p.wNum = i
		go pWorker(p, i)
	}
	return p.wg, nil
}

func (p *Parser) Stop() {
	close(p.out)
	fmt.Println("parser stopped")
}

func pWorker(p *Parser, i int) {
	fmt.Println(fmt.Sprintf("parse worker %d started", i))
	for msg := range p.in {
		fmt.Println("parser in", msg, len(p.in), msg.Shutdown, msg.Depth > p.maxDepth)
		//startTime := time.Now()
		if msg.Shutdown {
			break
		}
		if msg.Depth > p.maxDepth {
			break
		}
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
		p.out <- msg
		//fmt.Println(fmt.Sprintf("parse worker %d parsed url %s, %s", i, msg.Url, time.Now().Sub(startTime)))
		//fmt.Println("result", urlList)
	}
	p.wg.Done()
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
