package parser_test

import (
	"strconv"
	"testing"

	"github.com/relnod/calcgo/lexer"
	"github.com/relnod/calcgo/parser"
)

type Benchmark struct {
	name string
	str  string
}

var benchmarks []Benchmark

func init() {
	str := "1"
	appendStr := " + (1 * 1)"
	for i := 0; i < 20; i++ {
		s := str
		for j := 0; j < i*i; j++ {
			s += appendStr
		}
		benchmarks = append(benchmarks, Benchmark{
			name: strconv.Itoa(i * i),
			str:  s,
		})
	}
}

func BenchmarkParser(b *testing.B) {
	for _, bm := range benchmarks {
		b.Run("NonParallel"+bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				l := lexer.NewLexerFromString(bm.str)

				parser.ParseFromReader(l)
			}
		})

		b.Run("Parallel"+bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				l := lexer.NewBufferedLexerFromString(bm.str)

				parser.ParseFromReader(l)
			}
		})
	}
}
