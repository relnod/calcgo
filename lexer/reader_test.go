package lexer_test

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/relnod/calcgo/lexer"
)

type generator func(str string) lexer.BufferedReader

var _ = Describe("StaticBufferedReader", func() {
	describeBufferedReader(func(str string) lexer.BufferedReader {
		return lexer.NewBufferedReaderFromString(str)
	})
})

var _ = Describe("BufferedIOReader", func() {
	describeBufferedReader(func(str string) lexer.BufferedReader {
		r := bytes.NewReader([]byte(str))
		return lexer.NewBufferedReader(r)
	})
})

func describeBufferedReader(newReader generator) {
	It("has correct inital state", func() {
		r := lexer.NewBufferedReaderFromString("test")
		Expect(r.StartPos()).To(BeZero())
		Expect(r.CurrPos()).To(BeZero())
		Expect(r.Current()).To(BeZero())
		Expect(r.All()).To(BeEmpty())
	})

	It("returns correct next byte and has correct state afterwards", func() {
		r := lexer.NewBufferedReaderFromString("test")
		b, ok := r.Next()
		Expect(b).To(Equal(uint8('t')))
		Expect(ok).To(BeTrue())

		Expect(r.StartPos()).To(BeZero())
		Expect(r.CurrPos()).To(Equal(1))
		Expect(r.Current()).To(Equal(uint8('t')))
		Expect(r.All()).To(Equal([]byte("t")))
	})

	It("returns false at end of input", func() {
		r := lexer.NewBufferedReaderFromString("t")
		r.Next()
		b, ok := r.Next()
		Expect(b).To(BeZero())
		Expect(ok).To(BeFalse())
	})

	It("backups correctly", func() {
		r := lexer.NewBufferedReaderFromString("test")
		r.Next()
		r.Backup()

		Expect(r.CurrPos()).To(BeZero())
	})

	It("has correct state after resetting", func() {
		r := lexer.NewBufferedReaderFromString("test")
		r.Next()
		r.Reset()

		Expect(r.StartPos()).To(Equal(1))
		Expect(r.CurrPos()).To(Equal(1))
		Expect(r.All()).To(BeEmpty())

		r.Next()

		Expect(r.StartPos()).To(Equal(1))
		Expect(r.CurrPos()).To(Equal(2))
	})

	It("has correct state after backup and resetting", func() {
		r := lexer.NewBufferedReaderFromString("test")
		r.Next()
		r.Backup()
		r.Reset()

		Expect(r.StartPos()).To(Equal(0))
		Expect(r.CurrPos()).To(Equal(0))
		Expect(r.All()).To(BeEmpty())
	})
}
