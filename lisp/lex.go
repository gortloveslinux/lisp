package lisp

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"unicode"
)

type lexer struct {
	rr            io.RuneScanner
	err           error
	cl            bytes.Buffer // current line
	current, last rune         // last read rune
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
	tokenAtom
	tokenNumber
)

type token struct {
	typ TokenTyp
	val interface{}
}

func newLexer(rr io.RuneScanner) *lexer {
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
		case unicode.IsLetter(r):
			return l.readAtom()
		case unicode.IsNumber(r):
			return l.readNumber()
		case unicode.IsSpace(r):
			continue
		default:
			return &token{typ: tokenError, val: fmt.Sprintf("Unexpected token [%s]", string(r))}
		}
	}
}

func (l *lexer) read() (rune, error) {
	var err error
	l.last = l.current
	l.current, _, err = l.rr.ReadRune()
	if err != nil {
		return l.current, err
	}
	return l.current, nil
}

func (l *lexer) unread() error {
	err := l.rr.UnreadRune()
	if err != nil {
		return err
	}
	l.current = l.last
	l.last = 0
	return nil
}

// Comments ;.*\n
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

// Atoms [A-Za-z][A-Za-z0-9-_]*[A-Za-z]
func (l *lexer) readAtom() *token {
	var b bytes.Buffer
	b.WriteRune(l.current)
	for {
		r, err := l.read()
		switch {
		case err == io.EOF:
			if unicode.IsLetter(l.last) || unicode.IsNumber(l.last) {
				return &token{typ: tokenAtom, val: b.String()}
			} else {
				return &token{typ: tokenError, val: fmt.Sprintf("Invalid Atom [%s]", b.String())}
			}
		case unicode.IsLetter(r), unicode.IsDigit(r), r == '-', r == '_':
			b.WriteRune(r)
		default:
			err := l.unread()
			if err != nil {
				l.err = err
				return &token{typ: tokenError}
			}
			if unicode.IsLetter(l.current) || unicode.IsNumber(l.current) {
				return &token{typ: tokenAtom, val: b.String()}
			} else {
				return &token{typ: tokenError, val: fmt.Sprintf("Invalid Atom [%s]", b.String())}
			}
		}
	}
}

func (l *lexer) readNumber() *token {
	var b bytes.Buffer
	b.WriteRune(l.current)
	for {
		r, err := l.read()
		switch {
		case err == io.EOF:
			if unicode.IsNumber(l.last) {
				n, err := strconv.Atoi(b.String())
				if err != nil {
					return &token{typ: tokenError, val: fmt.Sprintf("Invalid Number [%s]: %s", b.String(), err.Error())}
				}
				return &token{typ: tokenNumber, val: n}
			}
		case unicode.IsNumber(r):
			b.WriteRune(r)
		default:
			err := l.unread()
			if err != nil {
				l.err = err
				return &token{typ: tokenError}
			}
			if unicode.IsNumber(l.current) {
				n, err := strconv.Atoi(b.String())
				if err != nil {
					return &token{typ: tokenError, val: fmt.Sprintf("Invalid Number [%s]: %s", b.String(), err.Error())}
				}
				return &token{typ: tokenNumber, val: n}
			} else {
				return &token{typ: tokenError, val: fmt.Sprintf("Invalid Number [%s]: %s", b.String(), err.Error())}
			}
		}

	}
}

func (t *token) String() string {
	return fmt.Sprintf("[%s:%v(%T)]", t.typ, t.val, t.val)
}
