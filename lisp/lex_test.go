package lisp

import (
	"fmt"
	"strings"
	"testing"
)

func TestComment(t *testing.T) {
	l := `; This is a comment
; This too is a comment
;;;;
`
	expected := []*token{
		&token{typ: tokenComment, txt: " This is a comment"},
		&token{typ: tokenComment, txt: " This too is a comment"},
		&token{typ: tokenComment, txt: ";;;"},
		&token{typ: tokenEOF},
	}
	sr := strings.NewReader(l)
	lxr := newLexer(sr)
	var tks []*token
	getTokens(lxr, &tks)
	if cmpTokenSlice(expected, tks) == false {
		t.Errorf("For test string\n%s\nExpected:\t%v\nGot:\t\t\t\t%v\n", l, expected, tks)
	}
}

func TestParen(t *testing.T) {
	tests := []struct {
		test     string
		expected []*token
	}{
		{`((()))`, []*token{
			&token{typ: tokenLParen},
			&token{typ: tokenLParen},
			&token{typ: tokenLParen},
			&token{typ: tokenRParen},
			&token{typ: tokenRParen},
			&token{typ: tokenRParen},
			&token{typ: tokenEOF}},
		},
		{`(;This is my opening comment
			))))`, []*token{
			&token{typ: tokenLParen},
			&token{typ: tokenComment, txt: "This is my opening comment"},
			&token{typ: tokenRParen},
			&token{typ: tokenRParen},
			&token{typ: tokenRParen},
			&token{typ: tokenRParen},
			&token{typ: tokenEOF}},
		},
	}

	for _, tst := range tests {
		sr := strings.NewReader(tst.test)
		lxr := newLexer(sr)
		var tks []*token
		getTokens(lxr, &tks)
		if cmpTokenSlice(tst.expected, tks) == false {
			t.Errorf("For test string %s\nExpected:\t%v\nGot:\t\t\t\t%v\n", tst.test, tst.expected, tks)
		}
	}
}

func getTokens(l *lexer, tks *[]*token) {
	t := &token{}
	for t.typ != tokenEOF {
		t = l.next()
		*tks = append(*tks, t)
	}
	fmt.Printf("tokens %v\n", tks)
}

func cmpTokenSlice(a []*token, b []*token) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v.typ != b[i].typ || v.txt != b[i].txt {
			return false
		}
	}
	return true
}
