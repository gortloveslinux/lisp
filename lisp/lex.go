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
	row        int
	col        int
	tcol       int
	trow       int
}

var EOFRUNE = rune(-1)
var ERRRUNE = rune(-2)

//go:generate go run golang.org/x/tools/cmd/stringer -type=tokenTyp
type tokenTyp int

const (
	tokenError tokenTyp = iota
	tokenEOF
	tokenComment
	tokenLParen
	tokenRParen
	tokenQuote
	tokenAtom
	tokenNumber
)

type token struct {
	typ tokenTyp
	val interface{}
	raw string
	row int
	col int
	err string
}

func newLexer(rr io.RuneScanner) *lexer {
	l := log.New(os.Stderr, "", 0)
	lex := &lexer{rr: rr, log: l, row: 1}
	return lex
}

func (l *lexer) next() *token {
	for {
		r := l.peek()
		l.tcol = l.col
		l.trow = l.row
		switch {
		case r == EOFRUNE:
			_ = l.read()
			return l.makeToken(tokenEOF, nil, "", "")
		case r == ERRRUNE:
			_ = l.read()
			return l.makeToken(tokenError, nil, "", "Rune Error")
		case r == ';':
			return l.readComment()
		case r == '\n', unicode.IsSpace(r):
			_ = l.read()
			continue
		case r == '(':
			_ = l.read()
			return l.makeToken(tokenLParen, nil, "(", "")
		case r == ')':
			_ = l.read()
			return l.makeToken(tokenRParen, nil, ")", "")
		case r == '\'':
			_ = l.read()
			return l.makeToken(tokenQuote, nil, "'", "")
		case unicode.IsLetter(r):
			return l.readAtom()
		case unicode.IsNumber(r), r == '.', r == '-':
			return l.readNumber()
		default:
			_ = l.read()
			return l.makeToken(tokenError, nil, "", fmt.Sprintf("Unexpected token[%s]", string(r)))
		}
	}
}

// Comments ;.*\n
func (l *lexer) readComment() *token {
	var b bytes.Buffer
	for {
		r := l.peek()
		switch r {
		case EOFRUNE, '\n':
			return l.makeToken(tokenComment, nil, b.String(), "")
		case ERRRUNE:
			return l.makeToken(tokenError, nil, "", "Rune Error")
		default:
			b.WriteRune(r)
			_ = l.read()
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
				return l.makeToken(tokenAtom, b.String(), b.String(), "")
			}
			return l.makeToken(tokenError, nil, b.String(), fmt.Sprintf("Invalid Atom[%s]", b.String()))
		}
	}
}

func (l *lexer) readNumber() *token {
	var b bytes.Buffer
	r := l.read()
	dec := r == '.'
	b.WriteRune(r)
	for {
		r := l.peek()
		switch {
		case unicode.IsNumber(r), r == '.' && !dec:
			dec = r == '.' || dec
			_ = l.read()
			b.WriteRune(r)
		default:
			if r == '\n' || r == '(' || r == ')' || r == ' ' {
				var (
					n   interface{}
					err error
				)
				if dec {
					n, err = strconv.ParseFloat(b.String(), 64)
				} else {
					n, err = strconv.Atoi(b.String())
				}
				if err != nil {
					return l.makeToken(tokenError, nil, b.String(), fmt.Sprintf("Invalid Number [%s]: %s", b.String(), err.Error()))
				}
				return l.makeToken(tokenNumber, n, b.String(), "")
			}
			return l.makeToken(tokenError, nil, b.String(), fmt.Sprintf("Invalid Number [%s]", b.String()))
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
	if r == '\n' {
		l.row = l.row + 1
		l.col = 0
	} else {
		l.col = l.col + 1
	}
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

func (l *lexer) makeToken(t tokenTyp, v interface{}, r, e string) *token {
	tok := &token{typ: t, val: v, raw: r, row: l.trow, col: l.tcol, err: e}
	return tok
}

func (l *lexer) peek() rune {
	if l.peeking {
		return l.curr
	}
	l.peeking = true
	r, _, err := l.rr.ReadRune()
	if r == '\n' {
		l.row = l.row + 1
		l.col = 0
	} else {
		l.col = l.col + 1
	}
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
	return fmt.Sprintf("[@(%d,%d)%s<%T>:%v,%s]", t.row, t.col, t.typ, t.val, t.val, t.raw)
}
