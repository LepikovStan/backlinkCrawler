// Package parser provides basic functions to parse html
package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"strings"
)

var count int

// Parser is the type, contains the basic Parse method
type Parser struct {
	Num int
}

// Parse parsing html and find all links on the page
// return list of finded links or error
func (p *Parser) Parse(body io.Reader) ([]string, error) {
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

// New function initialize new Parser instance and return pointer to it
func New() *Parser {
	p := new(Parser)
	count++
	p.Num = count
	return p
}
