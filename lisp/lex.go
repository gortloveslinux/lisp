package lisp

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"unicode"
)

type lexer struct {
	rr         io.RuneScanner
	curr, last rune // last read rune
	peeking    bool
	log        *log.Logger
}

var EOFRUNE = rune(-1)
var ERRRUNE = rune(-2)

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
	l := log.New(os.Stderr, "", 0)
	lex := &lexer{rr: rr, log: l}
	return lex
}

func (l *lexer) next() *token {
	for {
		r := l.peek()
		switch {
		case r == EOFRUNE:
			return &token{typ: tokenEOF}
		case r == ERRRUNE:
			return &token{typ: tokenError}
		case r == ';':
			return l.readComment()
		case r == '\n', unicode.IsSpace(r):
			_ = l.read()
			continue
		case r == '(':
			_ = l.read()
			return &token{typ: tokenLParen}
		case r == ')':
			_ = l.read()
			return &token{typ: tokenRParen}
		case r == '\'':
			_ = l.read()
			return &token{typ: tokenQuote}
		case unicode.IsLetter(r):
			return l.readAtom()
		case unicode.IsNumber(r):
			return l.readNumber()
		default:
			_ = l.read()
			return &token{typ: tokenError, val: fmt.Sprintf("Unexpected token[%s]", string(r))}
		}
	}
}

// Comments ;.*\n
func (l *lexer) readComment() *token {
	var b bytes.Buffer
	for {
		r := l.read()
		switch r {
		case EOFRUNE, '\n':
			return &token{tokenComment, b.String()}
		case ERRRUNE:
			return &token{typ: tokenError}
		default:
			b.WriteRune(r)
		}
	}
}

// Atoms [A-Za-z][A-Za-z0-9-_]*[A-Za-z0-9]
func (l *lexer) readAtom() *token {
	var b bytes.Buffer
	r := l.read()
	b.WriteRune(r)
	for {
		r := l.peek()
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r), r == '-', r == '_':
			_ = l.read()
			b.WriteRune(r)
		default:
			if unicode.IsLetter(l.last) || unicode.IsNumber(l.last) {
				return &token{typ: tokenAtom, val: b.String()}
			}
			return &token{typ: tokenError, val: fmt.Sprintf("Invalid Atom[%s]", b.String())}
		}
	}
}

func (l *lexer) readNumber() *token {
	var b bytes.Buffer
	r := l.read()
	b.WriteRune(r)
	for {
		r := l.peek()
		switch {
		case unicode.IsNumber(r):
			_ = l.read()
			b.WriteRune(r)
		default:
			if r == '\n' || r == '(' || r == ')' || r == ' ' {
				n, err := strconv.Atoi(b.String())
				if err != nil {
					return &token{typ: tokenError, val: fmt.Sprintf("Invalid Number [%s]: %s", b.String(), err.Error())}
				}
				return &token{typ: tokenNumber, val: n}
			}
			return &token{typ: tokenError, val: fmt.Sprintf("Invalid Number [%s]", b.String())}
		}
	}
}

func (l *lexer) read() rune {
	l.last = l.curr
	if l.peeking {
		l.peeking = false
		return l.curr
	}
	r, _, err := l.rr.ReadRune()
	if err != nil {
		if err == io.EOF {
			l.curr = EOFRUNE
			return EOFRUNE
		}
		l.log.Printf("Error reading rune: %s", err)
		l.curr = ERRRUNE
		return ERRRUNE
	}
	l.curr = r
	return r
}

func (l *lexer) peek() rune {
	if l.peeking {
		return l.curr
	}
	l.peeking = true
	r, _, err := l.rr.ReadRune()
	if err != nil {
		if err == io.EOF {
			l.curr = EOFRUNE
			return EOFRUNE
		}
		l.log.Printf("Error reading rune: %s", err)
		l.curr = ERRRUNE
		return ERRRUNE
	}
	l.curr = r
	return r
}

func (t *token) String() string {
	return fmt.Sprintf("[%s<%T>:%v]", t.typ, t.val, t.val)
}
