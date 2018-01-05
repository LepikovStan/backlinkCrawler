package crawler

import (
	"testing"
	"reflect"
	"log"
)

func TestGetStartList(t *testing.T) {
	result := GetStartList("input.txt")

	if reflect.TypeOf(result[0]).Kind() != reflect.String {
		log.Fatalf("Type of returned value is %s, expected string", reflect.TypeOf(result[0]))
	}
}

func BenchmarkGetStartList(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetStartList("input.txt")
	}
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New()
	}
}