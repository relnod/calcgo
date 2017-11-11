# Calcgo

[![Godoc](https://godoc.org/github.com/relnod/calcgo?status.svg)](https://godoc.org/github.com/relnod/calcgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/relnod/calcgo)](https://goreportcard.com/report/github.com/relnod/calcgo)

This is an experimental learning project, to try and better understand how a lexer and parser works.

## Syntax

```
a operator b
```

where a and b are either a number or another expression.

Valid operators are '```+```', '```-```', '```*```' and '```/```'.

Numbers can be either in the form of an integer (e.g '```12345```') or a floating point number (e.g. '```12.34```').

All calculations follow the "point before line rule". To break this rule brackets can be used.
```
1 + 2 * 3 = 7
```
but
```
(1 + 2) * 3 = 9
```

### Grammar

This gives the following grammar:

S -> SoS

S -> lSr

S -> n

#### Terminals
```n``` => ``n``umber ∈ { integer, float }

```o``` => ``o``perator ∈ { +, -, *, / }

```l``` => ``l``eft bracket ∈ { ( }

```r``` => ``r``ight bracket ∈ { ) }
