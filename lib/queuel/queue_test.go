package queuel

import (
	"log"
	"testing"
)

func TestTransformUrlToBacklink(t *testing.T) {
	urlsList := []string{"https://golang.org", "http://play.golang.org"}
	depth := 1

	result := TransformUrlToBacklink(urlsList, depth)

	if len(result) != 2 {
		log.Fatalf("Number of created Backlinks is %d, expected %d", len(result), 2)
	}

	if result[0].Url != urlsList[0] {
		log.Fatalf("Url of created element is '%s', expected '%s'", result[0].Url, urlsList[0])
	}

	if result[0].Depth != depth {
		log.Fatalf("Depth of created element is %d, expected %d", result[0].Depth, depth)
	}
}

func BenchmarkTransformUrlToBacklink(b *testing.B) {
	urlsList := []string{"https://golang.org", "http://play.golang.org"}

	for i := 0; i < b.N; i++ {
		TransformUrlToBacklink(urlsList, 1)
	}
}

func TestNew(t *testing.T) {
	result := New(1)
	_ = result
	//switch t := result.(type) {
	//case *Q:
	//
	//default:
	//	log.Fatalf("Type of returned valed is %s, expected *Q", t)
	//}
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New(1)
	}
}

func TestQ_GetBuffer(t *testing.T) {
	queue := New(1)
	urlsList := []string{"https://golang.org", "http://play.golang.org"}
	depth := 1
	ssl := TransformUrlToBacklink(urlsList, depth)
	queue.SetBuffer(ssl)
	result := queue.GetBuffer()

	if len(result) != 2 {
		log.Fatalf("Number of created Backlinks is %d, expected %d", len(result), 2)
	}

	if result[0].Url != urlsList[0] {
		log.Fatalf("Url of created element is '%s', expected '%s'", result[0].Url, urlsList[0])
	}

	if result[0].Depth != depth {
		log.Fatalf("Depth of created element is %d, expected %d", result[0].Depth, depth)
	}
}

func BenchmarkQ_GetBuffer(b *testing.B) {
	queue := New(1)
	urlsList := []string{"https://golang.org", "http://play.golang.org"}
	depth := 1
	ssl := TransformUrlToBacklink(urlsList, depth)
	queue.SetBuffer(ssl)

	for i := 0; i < b.N; i++ {
		queue.GetBuffer()
	}
}

func TestQ_PopBuffer(t *testing.T) {
	queue := New(1)
	urlsList := []string{"https://golang.org", "http://play.golang.org"}
	depth := 1
	ssl := TransformUrlToBacklink(urlsList, depth)
	queue.SetBuffer(ssl)
	resultLength := 1
	result := queue.PopBuffer(resultLength)

	if len(result) != resultLength {
		log.Fatalf("Number of created Backlinks is %d, expected %d", len(result), resultLength)
	}

	if result[0].Url != urlsList[0] {
		log.Fatalf("Url of created element is '%s', expected '%s'", result[0].Url, urlsList[0])
	}

	if result[0].Depth != depth {
		log.Fatalf("Depth of created element is %d, expected %d", result[0].Depth, depth)
	}
}

func BenchmarkQ_PopBuffer(b *testing.B) {
	queue := New(1)
	urlsList := []string{"https://golang.org", "http://play.golang.org"}
	depth := 1
	ssl := TransformUrlToBacklink(urlsList, depth)
	queue.SetBuffer(ssl)

	for i := 0; i < b.N; i++ {
		queue.PopBuffer(1)
	}
}

func TestQ_SetBuffer(t *testing.T) {
	queue := New(1)
	urlsList := []string{"https://golang.org", "http://play.golang.org"}
	depth := 1
	ssl := TransformUrlToBacklink(urlsList, depth)
	queue.SetBuffer(ssl)
	result := queue.GetBuffer()

	if len(result) != 2 {
		log.Fatalf("Number of created Backlinks is %d, expected %d", len(result), 2)
	}
}

func BenchmarkQ_SetBuffer(b *testing.B) {
	queue := New(1)
	urlsList := []string{"https://golang.org", "http://play.golang.org"}
	depth := 1
	ssl := TransformUrlToBacklink(urlsList, depth)

	for i := 0; i < b.N; i++ {
		queue.SetBuffer(ssl)
	}
}

func TestQ_Get(t *testing.T) {
	queue := New(1)
	urlsList := []string{"https://golang.org"}
	depth := 1
	ssl := TransformUrlToBacklink(urlsList, depth)
	queue.Set(ssl)

	result := queue.Get()

	if result.Url != urlsList[0] {
		log.Fatalf("Url of created element is '%s', expected '%s'", result.Url, urlsList[0])
	}

	if result.Depth != depth {
		log.Fatalf("Depth of created element is %d, expected %d", result.Depth, depth)
	}
}

func BenchmarkQ_Get(b *testing.B) {
	queue := New(1)
	urlsList := []string{"https://golang.org"}
	depth := 1
	ssl := TransformUrlToBacklink(urlsList, depth)

	for i := 0; i < b.N; i++ {
		queue.Set(ssl)
		queue.Get()
	}
}
