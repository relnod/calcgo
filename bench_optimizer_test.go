package calcgo_test

import (
	"testing"

	"github.com/relnod/calcgo"
)

var (
	str          = "(1 + 2) * 4 - (4 / 6) + a"
	calculations = 100
)

func BenchmarkInterpreterNoOptimizer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		interpreter := calcgo.NewInterpreter(str)
		for i := 0; i < calculations; i++ {
			interpreter.SetVar("a", 5.0)
			interpreter.GetResult()
		}
	}
}

func BenchmarkInterpreterOptimizer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		interpreter := calcgo.NewInterpreter(str)
		interpreter.EnableOptimizer()
		for i := 0; i < calculations; i++ {
			interpreter.SetVar("a", 5.0)
			interpreter.GetResult()
		}
	}
}
