# Calcgo

[![Build Status](https://travis-ci.org/relnod/calcgo.svg?branch=master)](https://travis-ci.org/relnod/calcgo)
[![codecov](https://codecov.io/gh/relnod/calcgo/branch/master/graph/badge.svg)](https://codecov.io/gh/relnod/calcgo)
[![Godoc](https://godoc.org/github.com/relnod/calcgo?status.svg)](https://godoc.org/github.com/relnod/calcgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/relnod/calcgo)](https://goreportcard.com/report/github.com/relnod/calcgo)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/faa47f159bbb47a497f602b6e8037c0d)](https://www.codacy.com/app/relnod/calcgo?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=relnod/calcgo&amp;utm_campaign=Badge_Grade)

This is an experimental learning project, to better understand the process of
lexing and parsing.

## Description
Calcgo exposes a lexer, parser and interpreter to get tokens, an ast and the
result of a basic mathematical calculation. All three functions accept a
language L(G) defined [here](#grammar).

The calculations follow basic math rules, like "multiplication and division
frist, then addition and subtraction" rule. To break this rule it is possible
to use brackets.
There needs to be at least one whitespace character between an operator an a
number. All other whitespace character get ignored by the lexer.


#### Lexer:
``` go
calcgo.Lex("(1 + 2) * 3")
```
#### Parser:
``` go
calcgo.Parse("(1 + 2) * 3")
```
#### Interpreter:
``` go
calcgo.Interpret("1 + 2 * 3")   // Result: 7
calcgo.Interpret("(1 + 2) * 3") // Result: 9
```

#### Interpreter with variable:
Calcgo supports variables. An instantiation of all variables have to be supplied
before interpreting.
```go
interpreter := calcgo.NewInterpreter("1 + a")
interpreter.SetVar("a", 1.0)
interpreter.GetResult() // Result: 2
interpreter.SetVar("a", 2.0)
interpreter.GetResult() // Result: 3
```

## Example
``` go
package main

import (
	"fmt"

	"github.com/relnod/calcgo"
)

func main() {
	number, _ := calcgo.Interpret("1 + 1")

	fmt.Println(number)
}
```

## Grammar

### Lexer:

|         | in       |          | out      |
|---------|----------|----------|----------|
| START   | 0        | NUMBER   |          |
| START   | 1-9      | INT      |          |
| START   | (        | START    | TLPAREN  |
| START   | )        | START    | TRPAREN  |
| START   | +        | START    | TOpPlus  |
| START   | -        | START    | TOpMinus |
| START   | *        | START    | TOpMult  |
| START   | /        | START    | TOpDiv   |
| START   | %        | START    | TOpMod   |
| START   | |        | START    | TOpOr    |
| START   | ^        | START    | TOpXor   |
| START   | &        | START    | TOpAnd   |
| START   | s        | S        |          |
| START   | c        | C        |          |
| START   | t        | T        |          |
| START   | a-z      | VAR      |          |
|         |          |          |          |
| NUMBER  | b        | BYTE     |          |
| NUMBER  | x        | HEX      |          |
| NUMBER  | .        | DEC      |          |
|         |          |          |          |
| BIN     | 0,1      | BIN2     |          |
|         |          |          |          |
| BIN2    | 0,1      | BIN2     |          |
| BIN2    | " "      | START    | TBin     |
|         |          |          |          |
| HEX     | 1-9, A-F | HEX2     |          |
|         |          |          |          |
| HEX2    | 1-0, A-F | HEX2     |          |
| HEX2    | " "      | START    | THex     |
|         |          |          |          |
| INT     | 1-9      | INT      |          |
| INT     | .        | DEC      |          |
| INT     | ^        | EXP      |          |
| INT     | " "      | START    | TInt     |
|         |          |          |          |
| DEC     | 1-9      | DEC      |          |
| DEC     | " "      | START    | TDec     |
|         |          |          |          |
| Exp     | 1-9      | Exp      |          |
| Exp     | " "      | START    | TExp     |
|         |          |          |          |
| S       | q        | SQ       |          |
| S       | i        | SI       |          |
| S       | a-z      | VAR      |          |
| S       | " "      | STAR     | TVar     |
|         |          |          |          |
| SQ      | r        | SQR      |          |
| SQ      | a-z      | VAR      |          |
| SQ      | " "      | STAR     | TVar     |
|         |          |          |          |
| SQR     | t        | SQRT     |          |
| SQR     | a-z      | VAR      |          |
| SQR     | " "      | STAR     | TVar     |
|         |          |          |          |
| SQRT    | (        | START    | TFnSqrt  |
| SQRT    | a-z      | VAR      |          |
| SQRT    | " "      | STAR     | TVar     |
|         |          |          |          |
| SI      | n        | SIN      |          |
| SI      | a-z      | VAR      |          |
| SI      | " "      | STAR     | TVar     |
|         |          |          |          |
| SIN     | (        | START    | TFnSin   |
| SIN     | a-z      | VAR      |          |
| SIN     | " "      | STAR     | TVar     |
|         |          |          |          |
| C       | o        | CO       |          |
| C       | a-z      | VAR      |          |
| C       | " "      | STAR     | TVar     |
|         |          |          |          |
| CO      | s        | COS      |          |
| CO      | a-z      | VAR      |          |
| CO      | " "      | STAR     | TVar     |
|         |          |          |          |
| Cos     | (        | START    | TFnCos   |
| COS     | a-z      | VAR      |          |
| COS     | " "      | STAR     | TVar     |
|         |          |          |          |
| T       | a        | TA       |          |
| T       | a-z      | VAR      |          |
| T       | " "      | STAR     | TVar     |
|         |          |          |          |
| TA      | n        | TAN      |          |
| TA      | a-z      | VAR      |          |
| TA      | " "      | STAR     | TVar     |
|         |          |          |          |
| TAN     | (        | START    | TFnTan   |
| TAN     | a-z      | VAR      |          |
| TAN     | " "      | STAR     | TVar     |
|         |          |          |          |
| VAR     | a-z      | VAR      |          |
| VAR     | " "      | START    | TVar     |


### Lexer:
The lexer accepts the language L(G), where G is the following grammar:

G = (N,T,P,S)

N = {S}

T = {n, o, l, r}

P contains the following rules:

S → SS | e | nN | n | cC | c | o | l | r

N → nN | s

C → cC | s

#### Terminals
n ∈ { '0', '1', '2', '3', '4', '5', '6', '7', '8', '9' }

c ∈ { 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z' }

d ∈  { '.' }

o ∈ { '+', '-', '*'*, '/' }

l ∈ { '(' }

r ∈ { ')' }

s ∈  { ' ' }

e ∈  { EOF }

### Parser:
The parser accepts the language L(G), where G is the following grammar:

G = (N,T,P,S)

N = {S}

T = {n, o, l, r}

P contains the following rules:

S → SoS | lSr | fSr | n

#### Terminals
n ∈ { TInteger, TDecimal, TVariable }

o ∈ { TOperatorPlus, TOperatorMinus, TOperatorMult, TOperatorDiv }

f ∈ { TFuncSqrt }

l ∈ { TLeftBracket }

r ∈ { TRightBracket }

## Tests and Benchmarks

#### Running Tests

Run tests with ```go test -v -race``` or with ```goconvey``` to see live test
results in a browser.

#### Running Benchmarks

Benchmarks can be tun with ```go test -run=^$ -bench=.```

To see the differences between current branch and master run ```./scripts/benchcmp.sh -n 5```

## License

This project is licensed under the MIT License. See the
[LICENSE](../master/LICENSE) file for details.
