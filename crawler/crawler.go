// Package crawler provides basic functions to crawl html web pages
package crawler

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var count int

func readFile(path string) []string {
	var result []string
	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return result
}

// GetStartList is function for receiving start list of web pages to crawl
func GetStartList() []string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return readFile(fmt.Sprintf("%s/%s", dir, "crawler/input.txt"))
	//return []string{
	//	// "https://yandex.ru/yandsearch?text=golang%20error%20type&lr=2",
	//	//"https://gobyexample.com/errors",
	//	// "https://golang.org/pkg/sync",
	//	//"https://developers.google.com/products/",
	//	//"http://www.tattyworld.net",
	//	//"https://yandex.ru/",
	//	"https://golang.org/",
	//	"",
	//	"https://gobyexample.com/errors",
	//}
}

// Crawler is the type, contains the basic Crawl method
type Crawler struct {
	Num int
}

// Crawl crawl web page via http.Get and return 200 status code or
// error if exists
func (c *Crawler) Crawl(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		return resp.Body, nil
	}

	err = errors.New(fmt.Sprintf("Wrong response status: %d for %s", resp.StatusCode, url))
	return nil, err
}

// New function initialize new Crawler instance and return pointer to it
func New() *Crawler {
	c := new(Crawler)
	count++
	c.Num = count
	return c
}
