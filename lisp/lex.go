package lisp

import (
	"bytes"
	"fmt"
	"io"
)

type lexer struct {
	rr   io.RuneReader
	err  error
	cl   bytes.Buffer // current line
	last rune         // last read rune
}

//go:generate stringer -type=TokenTyp
type TokenTyp int

const (
	tokenError TokenTyp = iota
	tokenEOF
	tokenComment
	tokenLParen
	tokenRParen
	tokenQuote
)

type token struct {
	typ TokenTyp
	txt string
}

func newLexer(rr io.RuneReader) *lexer {
	return &lexer{rr: rr}
}

func (l *lexer) next() *token {
	for {
		r, err := l.read()
		if err != nil {
			if err == io.EOF {
				return &token{typ: tokenEOF}
			}
		}
		switch {
		case r == ';':
			return l.readComment()
		case r == '\n':
			continue
		case r == '(':
			return &token{typ: tokenLParen}
		case r == ')':
			return &token{typ: tokenRParen}
		case r == '\'':
			return &token{typ: tokenQuote}
			//case unicode.IsLetter(r):
			//return l.readAtom()
		}
	}
}

func (l *lexer) read() (rune, error) {
	r, _, err := l.rr.ReadRune()
	if err != nil {
		return r, err
	}
	l.last = r
	return r, nil
}

func (l *lexer) readComment() *token {
	var b bytes.Buffer
	for {
		r, err := l.read()
		if err != nil {
			if err == io.EOF {
				return &token{tokenComment, b.String()}
			}
			return &token{typ: tokenError}
		}
		if r == '\n' {
			return &token{tokenComment, b.String()}
		}
		b.WriteRune(r)
	}
}

/*
func (l *lexer) readAtom() *token {
	var b bytes.Buffer
	for {
		r, err := l.read()
		if err != nil {
			if err == io.EOF {
				return &token{typ: tokenEOF}
			}
			l.err = err
			return &token{typ: tokenError}
		}
	}
}
*/

func (t *token) String() string {
	return fmt.Sprintf("[%s:%s]", t.typ, t.txt)
}
