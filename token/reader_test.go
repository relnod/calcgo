package token_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/relnod/calcgo/token"
)

var _ = Describe("StaticReader", func() {
	It("reads all tokens", func() {
		tokens := []token.Token{
			{Type: token.Int, Value: "test"},
			{Type: token.Dec, Value: "test2"},
			{Type: token.And, Value: "test3"},
			{Type: token.ParenL, Value: "test4"},
		}
		r := token.NewStaticReader(tokens)

		tokensRead := make([]token.Token, 0, len(tokens))
		for {
			t := r.Read()
			if t.Type == token.EOF {
				break
			}

			tokensRead = append(tokensRead, t)
		}

		Expect(tokensRead).To(Equal(tokens))
	})
})
