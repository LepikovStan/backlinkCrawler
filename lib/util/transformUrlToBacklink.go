package util

import "github.com/LepikovStan/backlinkCrawler/lib/queue"

func TransformUrlToBacklink(urls []string, depth int) []queue.Backlink {
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
