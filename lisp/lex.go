package lisp

import (
	"bytes"
	"io"
	"unicode"
)

type lexer struct {
	rr  io.RuneReader
	ln  int // line number
	cn  int // column number
	err error
	cl  bytes.Buffer // current line
}

type TokenTyp int

const (
	tokenError TokenTyp = iota
	tokenEOF
	tokenComment
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
				// TODO return tokenEOF
				return &token{typ: tokenEOF}
			}
		}
		switch {
		case r == ';':
			return l.readComment()
		case r == '\n':
		case r == '(':
		case r == ')':
		case r == '\'':
		case unicode.IsLetter(r):
		}
	}
}

func (l *lexer) read() (rune, error) {
	r, _, err := l.rr.ReadRune()
	if err != nil {
		return r, err
	}
	l.recordRead(r)
	return r, nil
}

func (l *lexer) readComment() *token {
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
		b.WriteRune(r)
		if r == '\n' {
			l.markNewline()
			return &token{tokenComment, b.String()}
		}
	}
}

func (l *lexer) markNewline() {
	l.cn = 0
	l.ln++
	l.cl.Reset()
}

func (l *lexer) recordRead(r rune) {
	l.cn++
	l.cl.WriteRune(r)
}
