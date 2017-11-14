package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"strings"
)

type Parser struct{}

func (p *Parser) Parse(body io.Reader) ([]string, error) {
	//fmt.Println("start Parse")
	doc, err := goquery.NewDocumentFromReader(body)
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

func New() *Parser {
	return new(Parser)
}
