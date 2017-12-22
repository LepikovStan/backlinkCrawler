// Package queue provides methods for runtime storages of
// urls to crawl
package queue

import (
	"io"
	"sync"
)

// Backlink type contains fields for crawl and parse backlinks
type Backlink struct {
	Url    string
	Body   io.Reader
	BLList []Backlink
	Error  error
	Depth  int
}

// Q is main struct in package to work with runtime storages
type Q struct {
	mux          sync.Mutex
	q            chan Backlink
	buffer       []Backlink
	popbuffermux sync.Mutex
	setbuffermux sync.Mutex
}

// GetBuffer func return elements of buffer
func (queue *Q) GetBuffer() []Backlink {
	return queue.buffer
}

// SetBuffer func set given elements to buffer
func (queue *Q) SetBuffer(n []Backlink) {
	queue.setbuffermux.Lock()
	defer queue.setbuffermux.Unlock()

	queue.buffer = append(queue.buffer, n...)
}

// PopBuffer func return n head elements of buffer and make buffer less to n elements
func (queue *Q) PopBuffer(n int) []Backlink {
	queue.popbuffermux.Lock()
	defer queue.popbuffermux.Unlock()

	if len(queue.buffer) < n {
		res := queue.buffer
		queue.buffer = []Backlink{}
		return res
	}

	if len(queue.buffer) > 0 {
		res := queue.buffer[:n]
		queue.buffer = queue.buffer[n:]
		return res
	}
	return nil
}

// Init func initialize inner channel and buffer slice
func (queue *Q) Init(bufferSize int) {
	queue.q = make(chan Backlink, bufferSize)
	queue.buffer = []Backlink{}
}

// Write func set to inner channel
func (queue *Q) Set(ssl []Backlink) {
	for _, s := range ssl {
		queue.q <- s
	}
}

// Reag func get from inner channel
func (queue *Q) Get() Backlink {
	return <-queue.q
}

// GetChan func return
func (queue *Q) GetChan() chan Backlink {
	return queue.q
}

// New function initialize new Q instance and return pointer to it
func New(bufferSize int) *Q {
	result := new(Q)
	result.Init(bufferSize)
	return result
}
