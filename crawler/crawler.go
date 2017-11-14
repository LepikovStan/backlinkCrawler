package crawler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Crawler struct{}

func (c *Crawler) Crawl(url string) (io.Reader, error) {
	//fmt.Println("start Crawl", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		return resp.Body, nil
	}

	err = errors.New(fmt.Sprintf("Wrong status: %d for %s", resp.StatusCode, url))
	return nil, err
}

func New() *Crawler {
	return new(Crawler)
}
