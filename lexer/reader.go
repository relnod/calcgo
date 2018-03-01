package lexer

import (
	"io"
)

// BufferedReader describes the interface for a buffered reader, that stores
// the read contents in a buffer until it gets reset.
type BufferedReader interface {
	// StartPos returns the position of the total input since the last reset.
	StartPos() int

	// CurrPos returns the current position relative to the total input.
	CurrPos() int

	// Current returns the byte at the current position.
	Current() byte

	// All returns the whole content of the buffer.
	All() []byte

	// Next returns the next byte read from the input. If an error occured
	// during the read process or if the reader is at the end of its input the
	// second return value is false.
	Next() (byte, bool)

	// Backup moves the position of the buffer one to the left.
	Backup()

	// Reset resets the buffer content.
	Reset()
}

// BufferedIOReader implements the BufferedReader with an io.Reader as input.
type BufferedIOReader struct {
	reader io.Reader
	buf    []byte
	pos    int
	start  int
}

// NewBufferedReader returns a new buffered reader, that reads from the given
// io.Reader.
func NewBufferedReader(reader io.Reader) *BufferedIOReader {
	return &BufferedIOReader{
		reader: reader,
		buf:    make([]byte, 0),
		pos:    0,
		start:  0,
	}
}

// StartPos returns the position of the total input since the last reset.
func (r *BufferedIOReader) StartPos() int {
	return r.start
}

// CurrPos returns the current position relative to the total input.
func (r *BufferedIOReader) CurrPos() int {
	return r.pos + r.start
}

// Current returns the byte at the current position.
func (r *BufferedIOReader) Current() byte {
	return r.buf[r.pos-1]
}

// All returns the whole content of the buffer.
func (r *BufferedIOReader) All() []byte {
	return r.buf[0:r.pos]
}

// Next returns the next byte read from the input. If an error occured
// during the read process or if the reader is at the end of its input the
// second return value is false.
func (r *BufferedIOReader) Next() (byte, bool) {
	if r.pos < len(r.buf) {
		r.pos++
		return r.Current(), true
	}

	bt := make([]byte, 1, 1)
	_, err := r.reader.Read(bt)
	if err != nil {
		return 0, false
	}

	r.buf = append(r.buf, bt[0])

	r.pos++
	return r.Current(), true
}

// Backup moves the position of the buffer one to the left.
func (r *BufferedIOReader) Backup() {
	r.pos--
}

// Reset resets the buffer content.
func (r *BufferedIOReader) Reset() {
	r.buf = r.buf[r.pos:len(r.buf)]
	r.start = r.CurrPos()
	r.pos = 0
}

// StaticBufferedReader implements a buffered reader with a string as input.
type StaticBufferedReader struct {
	str   string
	pos   int
	start int
}

// NewBufferedReaderFromString returns a new buffered reader, with the given
// string as input.
func NewBufferedReaderFromString(str string) *StaticBufferedReader {
	return &StaticBufferedReader{
		str:   str,
		pos:   0,
		start: 0,
	}
}

// StartPos returns the position of the total input since the last reset.
func (r *StaticBufferedReader) StartPos() int {
	return r.start
}

// CurrPos returns the current position relative to the total input.
func (r *StaticBufferedReader) CurrPos() int {
	return r.pos
}

// Current returns the byte at the current position.
func (r *StaticBufferedReader) Current() byte {
	if r.pos == 0 {
		return 0
	}
	return r.str[r.pos-1]
}

// All returns the whole content of the buffer.
func (r *StaticBufferedReader) All() []byte {
	return []byte(r.str[r.start:r.pos])
}

// Next returns the next byte read from the input. If an error occured
// during the read process or if the reader is at the end of its input the
// second return value is false.
func (r *StaticBufferedReader) Next() (byte, bool) {
	if r.pos >= len(r.str) {
		return 0, false
	}

	r.pos++
	return r.Current(), true
}

// Backup moves the position of the buffer one to the left.
func (r *StaticBufferedReader) Backup() {
	r.pos--
}

// Reset resets the buffer content.
func (r *StaticBufferedReader) Reset() {
	r.start = r.pos
}
