# Calcgo

[![Build Status](https://travis-ci.org/relnod/calcgo.svg?branch=master)](https://travis-ci.org/relnod/calcgo)
[![codecov](https://codecov.io/gh/relnod/calcgo/branch/master/graph/badge.svg)](https://codecov.io/gh/relnod/calcgo)
[![Godoc](https://godoc.org/github.com/relnod/calcgo?status.svg)](https://godoc.org/github.com/relnod/calcgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/relnod/calcgo)](https://goreportcard.com/report/github.com/relnod/calcgo)

This is an experimental learning project, to better understand the process of
lexing and parsing.

## Description
Calcgo exposes a lexer, parser and interpreter to get tokens, an ast and the
result of a basic mathematical calculation. All three functions accept a
language L(G) defined [here](#grammar).

The calculations follow basic math rules, like "multiplication and division
first, then addition and subtraction" rule. To break this rule it is possible
to use brackets.
There needs to be at least one whitespace character between an operator an a
number. All other whitespace character get ignored by the lexer.


#### Lexer:
``` go
lexer.Lex("(1 + 2) * 3")
```
#### Parser:
``` go
parser.Parse("(1 + 2) * 3")
```
#### Interpreter:
``` go
interpreter.Interpret("1 + 2 * 3")   // Result: 7
interpreter.Interpret("(1 + 2) * 3") // Result: 9
```

#### Interpreter with variable:
Calcgo supports variables. An instantiation of all variables has to be supplied
before interpreting.
```go
i := interpreter.NewInterpreter("1 + a")
i.SetVar("a", 1.0)
i.GetResult() // Result: 2
i.SetVar("a", 2.0)
i.GetResult() // Result: 3
```

## Example
``` go
package main

import (
	"fmt"

	"github.com/relnod/calcgo"
)

func main() {
	number, _ := calcgo.Calc("1 + 1")

	fmt.Println(number)
}
```

## Tests and Benchmarks

#### Running Tests

Run tests with ```go test -v -race ./...``` or with ```ginkgo -r -v```.

#### Running Benchmarks

Benchmarks can be tun with ```go test -run=^$ -bench=. ./...```

To see the differences between current branch and master run ```./scripts/benchcmp.sh -n 5```

## License

This project is licensed under the MIT License. See the
[LICENSE](../master/LICENSE) file for details.
