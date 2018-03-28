package main

import (
	"log"
	"testing"
)

func TestQueue_Push(t *testing.T) {
	q := NewOueue()
	testArr := []int{2, 4, 7, 9, 3, 5, 2}
	for i := 0; i < len(testArr); i++ {
		q.Push(testArr[i])
	}
	headItem := q.Pop()

	if headItem != testArr[0] {
		log.Fatalf("Expecting %d but got %d", testArr[0], headItem)
	}

	headItem = q.Pop()

	if headItem != testArr[1] {
		log.Fatalf("Expecting %d but got %d", testArr[1], headItem)
	}
}

func TestQueue_Pop(t *testing.T) {
	q := NewOueue()
	testArr := []int{2, 4, 7, 9, 3, 5, 2, 10, 12, 1, 23, 16, 100, 23, 3456, 0, 12, 32, 45, 15, 32}
	for i := 0; i < len(testArr); i++ {
		q.Push(testArr[i])
	}
	headItem := q.Pop()

	if headItem != testArr[0] {
		log.Fatalf("Expecting %d but got %d", testArr[0], headItem)
	}

	headItem = q.Pop()

	if headItem != testArr[1] {
		log.Fatalf("Expecting %d but got %d", testArr[1], headItem)
	}
}

func TestQueue_Push2(t *testing.T) {
	q := NewOueue()
	testArr := []int{2, 4, 7, 9, 3, 5, 2, 10, 12, 1}
	for i := 0; i < len(testArr); i++ {
		q.Push(testArr[i])
	}
	headItem := q.Pop()

	if headItem != testArr[0] {
		log.Fatalf("Expecting %d but got %d", testArr[0], headItem)
	}

	headItem = q.Pop()

	if headItem != testArr[1] {
		log.Fatalf("Expecting %d but got %d", testArr[1], headItem)
	}

	testArr2 := []int{15, 17}
	for i := 0; i < len(testArr2); i++ {
		q.Push(testArr2[i])
	}
}
