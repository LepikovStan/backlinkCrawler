package parser

import (
	"testing"
	"github.com/LepikovStan/backlinkCrawler/crawler"
	"log"
)

func TestParser_Parse(t *testing.T) {
	cr := crawler.New()
	pr := New()
	body, err := cr.Crawl("http://golang.org")

	if err != nil {
		log.Fatal(err)
	}

	result, err := pr.Parse(body)

	if err != nil {
		log.Fatal(err)
	}

	if len(result) != 4 {
		log.Fatal("Number of parsed links must be 4")
	}
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New()
	}
}

func BenchmarkParser_Parse(b *testing.B) {
	cr := crawler.New()
	pr := New()
	body, err := cr.Crawl("http://golang.org")

	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		pr.Parse(body)
	}
}
