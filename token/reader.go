package token

// Reader defines a token token reader interface.
type Reader interface {
	// Read returns the next token from the reader. Should return a token with
	// type EOF as last element.
	Read() Token
}

// StaticReader converts a list of tokens into a token reader.
type StaticReader struct {
	tokens []Token
	index  int
}

// NewStaticReader returns a new static token reader. Takes a list of tokens as
// input.
func NewStaticReader(tokens []Token) *StaticReader {
	return &StaticReader{tokens: tokens, index: 0}
}

// Read reads the next token from the token list. If it reached the end of the
// list returns a new token with type EOF.
func (r *StaticReader) Read() Token {
	if r.index >= len(r.tokens) {
		return Token{Type: EOF}
	}

	t := r.tokens[r.index]

	r.index++

	return t
}
