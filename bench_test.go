package calcgo_test

import (
	"testing"

	"github.com/relnod/calcgo"
	"github.com/relnod/calcgo/lexer"
	"github.com/relnod/calcgo/parser"
)

var (
	str1 = ""
	str2 = "(1 + 2) * 3"
	str3 = "(1 + 2) * 3 + (((2 + 1) * 3 / (5 - 1)) + 1 / 3) - 2 / 3"
)

func Benchmark1Lexer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		lexer.Lex(str1)
	}
}

func Benchmark2Lexer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		lexer.Lex(str2)
	}
}

func Benchmark3Lexer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		lexer.Lex(str3)
	}
}

func Benchmark1Parser(b *testing.B) {
	for n := 0; n < b.N; n++ {
		parser.Parse(str1)
	}
}

func Benchmark2Parser(b *testing.B) {
	for n := 0; n < b.N; n++ {
		parser.Parse(str2)
	}
}

func Benchmark3Parser(b *testing.B) {
	for n := 0; n < b.N; n++ {
		parser.Parse(str3)
	}
}

func Benchmark1Interpreter(b *testing.B) {
	for n := 0; n < b.N; n++ {
		calcgo.Interpret(str1)
	}
}

func Benchmark2Interpreter(b *testing.B) {
	for n := 0; n < b.N; n++ {
		calcgo.Interpret(str2)
	}
}

func Benchmark3Interpreter(b *testing.B) {
	for n := 0; n < b.N; n++ {
		calcgo.Interpret(str3)
	}
}
