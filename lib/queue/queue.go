package queue

import (
	"io"
	"sync"
)

type Backlink struct {
	Url    string
	Body   io.Reader
	BLList []Backlink
	Error  error
	Depth  int
}

type Q struct {
	mux sync.Mutex
	q   chan Backlink
	buffer []Backlink
	popbuffermux sync.Mutex
	setbuffermux sync.Mutex
}

func (queue *Q) SetBuffer(n []Backlink) {
	queue.setbuffermux.Lock()
	defer queue.setbuffermux.Unlock()

	queue.buffer = append(queue.buffer, n...)
}
func (queue *Q) PopBuffer(n int) []Backlink {
	queue.popbuffermux.Lock()
	defer queue.popbuffermux.Unlock()

	res := queue.buffer[:n]
	queue.buffer = queue.buffer[n:]
	return res
}
func (queue *Q) Init(bufferSize int) {
	queue.q = make(chan Backlink, bufferSize)
	queue.buffer = []Backlink{}
}

func (queue *Q) Write(ssl []Backlink) {
	for _, s := range ssl {
		queue.q <- s
	}
}

func (queue *Q) Read() Backlink {
	return <-queue.q
}

func (queue *Q) GetChan() chan Backlink {
	return queue.q
}

func New(bufferSize int) *Q {
	result := new(Q)
	result.Init(bufferSize)
	return result
}
