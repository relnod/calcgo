package lexer_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/relnod/calcgo/lexer"
	"github.com/relnod/calcgo/token"
)

var _ = Describe("BufferedLexer", func() {
	It("behaves the same as normal lexer", func() {
		l1 := lexer.NewLexerFromString("1 + 2")
		l2 := lexer.NewBufferedLexerFromString("1 + 2")

		for {
			t1 := l1.Read()
			t2 := l2.Read()
			Expect(t1).To(Equal(t2))

			if t1.Type == token.EOF {
				break
			}
		}
	})
})
