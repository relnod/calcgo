package calcgo

import "github.com/relnod/calcgo/interpreter"

// Calc calculates a numerical expression. May return any number of errors,
// that occur during lexing, parsing or interpreting.
func Calc(expression string) (float64, []error) {
	return interpreter.Interpret(expression)
}
