package calcgo_test

import (
	"testing"

	"github.com/relnod/calcgo"
)

func BenchmarkInterpreterNoVars(b *testing.B) {
	for n := 0; n < b.N; n++ {
		calcgo.Interpret("(1 + 2) * 4 - (4 / 6)")
		calcgo.Interpret("(2 + 2) * 4 - (4 / 6)")
		calcgo.Interpret("(3 + 2) * 4 - (4 / 6)")
		calcgo.Interpret("(4 + 2) * 4 - (4 / 6)")
		calcgo.Interpret("(5 + 2) * 4 - (4 / 6)")
	}
}

func BenchmarkInterpreterVars(b *testing.B) {
	for n := 0; n < b.N; n++ {
		i := calcgo.NewInterpreter("(a + 2) * 4 - (4 / 6)")
		i.SetVar("a", 1.0)
		i.GetResult()
		i.SetVar("a", 2.0)
		i.GetResult()
		i.SetVar("a", 3.0)
		i.GetResult()
		i.SetVar("a", 4.0)
		i.GetResult()
		i.SetVar("a", 5.0)
		i.GetResult()
	}
}

func BenchmarkInterpreterVars1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		i := calcgo.NewInterpreter("(a + 2) * 4 - (4 / 6)")
		i.SetVar("a", 1.0)
		i.GetResult()
	}
}

func BenchmarkInterpreterVars2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		i := calcgo.NewInterpreter("(a + 2) * 4 - (4 / 6)")
		i.SetVar("a", 1.0)
		i.GetResult()
		i.SetVar("a", 1.0)
		i.GetResult()
	}
}
