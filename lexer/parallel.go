package lexer

import (
	"io"

	"github.com/relnod/calcgo/token"
)

// BufferedLexer wraps the lexer around a buffer,
type BufferedLexer struct {
	buf chan token.Token
	l   *Lexer
}

// NewBufferedlLexer creates a new buffered lexer and runs the lexer in a go
// routine.
func NewBufferedlLexer(r io.Reader) *BufferedLexer {
	b := &BufferedLexer{
		buf: make(chan token.Token, 5),
		l:   NewLexer(r),
	}

	b.run()

	return b
}

// NewBufferedLexerFromString creates a new buffered lexer and runs the lexer
// in a go routine.
func NewBufferedLexerFromString(str string) *BufferedLexer {
	b := &BufferedLexer{
		buf: make(chan token.Token, len(str)/4),
		l:   NewLexerFromString(str),
	}

	go b.run()

	return b
}

// Read returns the next token from the token chanel.
func (b *BufferedLexer) Read() token.Token {
	return <-b.buf
}

// Start starts the parallel lexer
func (b *BufferedLexer) run() {
	for {
		t := b.l.Read()
		if t.Type == token.EOF {
			break
		}
		b.buf <- t
	}
	close(b.buf)
}
