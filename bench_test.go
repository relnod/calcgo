package calcgo_test

import (
	"testing"

	"github.com/relnod/calcgo"
)

var str = "(1 + 2) * 3 + (((2 + 1) * 3 / (5 - 1)) + 1 / 3) - 2 / 3"

func BenchmarkLexer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		calcgo.Lex(str)
	}
}

func BenchmarkParser(b *testing.B) {
	for n := 0; n < b.N; n++ {
		calcgo.Parse(str)
	}
}

func BenchmarkInterpreter(b *testing.B) {
	for n := 0; n < b.N; n++ {
		calcgo.Interpret(str)
	}
}
