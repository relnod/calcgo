package lexer_test

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/relnod/calcgo/lexer"
	"github.com/relnod/calcgo/token"
)

var _ = Describe("BufferedLexer", func() {
	str := "1 + 2"

	testBufferedLexer := func(l1 *lexer.BufferedLexer) {
		It("behaves the same as normal lexer", func() {
			l2 := lexer.NewBufferedLexerFromString(str)

			for {
				t1 := l1.Read()
				t2 := l2.Read()
				Expect(t1).To(Equal(t2))

				if t1.Type == token.EOF {
					break
				}
			}
		})
	}

	Describe("from io.Reader", func() {
		testBufferedLexer(lexer.NewBufferedlLexer(bytes.NewReader([]byte(str))))
	})

	Describe("from string", func() {
		testBufferedLexer(lexer.NewBufferedLexerFromString(str))
	})
})
