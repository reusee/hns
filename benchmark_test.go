package nw

import (
	"io"
	"os"
	"testing"
)

func getTestHtml() io.Reader {
	f, err := os.Open("test.html")
	if err != nil {
		panic(err)
	}
	return f
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Parse(getTestHtml())
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWalk(b *testing.B) {
	root, err := Parse(getTestHtml())
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		root.Walk(Css("a", func(*Node) {}))
	}
}
