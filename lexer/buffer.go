package lexer

import (
	"io"
)

type Buffer struct {
	reader io.Reader
	buf    []byte
	pos    int
	start  int
}

func newBuffer(reader io.Reader) Buffer {
	return Buffer{
		reader: reader,
		buf:    make([]byte, 0),
		pos:    0,
		start:  0,
	}
}

func (b *Buffer) reset() {
	b.buf = b.buf[b.pos:len(b.buf)]
	b.start += b.pos
	b.pos = 0
}

func (b *Buffer) StartPos() int {
	return b.start
}

func (b *Buffer) Pos() int {
	return b.pos + b.start
}

func (b *Buffer) current() byte {
	return b.buf[b.pos-1]
}

func (b *Buffer) all() []byte {
	return b.buf[0:b.pos]
}

func (b *Buffer) next() (byte, bool) {
	if b.pos < len(b.buf) {
		b.pos++
		return b.current(), true
	}

	bt := make([]byte, 1, 1)
	_, err := b.reader.Read(bt)
	if err != nil {
		return 0, false
	}

	b.buf = append(b.buf, bt[0])

	b.pos++
	return b.current(), true
}

func (b *Buffer) backup() {
	b.pos--
}
