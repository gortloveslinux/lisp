package lisp

import (
	"fmt"
	"strings"
	"testing"
)

func TestComment(t *testing.T) {
	tests := []struct {
		test     string
		expected []*token
	}{
		{`; This is a comment
; This too is a comment
;;;;
`, []*token{
			&token{typ: tokenComment, txt: " This is a comment"},
			&token{typ: tokenComment, txt: " This too is a comment"},
			&token{typ: tokenComment, txt: ";;;"},
			&token{typ: tokenEOF},
		}},
	}

	for _, tst := range tests {
		sr := strings.NewReader(tst.test)
		lxr := newLexer(sr)
		var tks []*token
		getTokens(lxr, &tks)
		if cmpTokenSlice(tst.expected, tks) == false {
			t.Errorf("For test string\n%s\nExpected:\t%v\nGot:\t\t\t\t%v\n", tst.test, tst.expected, tks)
		}
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

func TestAtom(t *testing.T) {
	tests := []struct {
		test     string
		expected []*token
	}{
		{`(test())`, []*token{
			&token{typ: tokenLParen},
			&token{typ: tokenAtom, txt: "test"},
			&token{typ: tokenLParen},
			&token{typ: tokenRParen},
			&token{typ: tokenRParen},
			&token{typ: tokenEOF}},
		},
		{`(test(test foo_roo)) ;This is a comment`, []*token{
			&token{typ: tokenLParen},
			&token{typ: tokenAtom, txt: "test"},
			&token{typ: tokenLParen},
			&token{typ: tokenAtom, txt: "test"},
			&token{typ: tokenAtom, txt: "foo_roo"},
			&token{typ: tokenRParen},
			&token{typ: tokenRParen},
			&token{typ: tokenComment, txt: "This is a comment"},
			&token{typ: tokenEOF}},
		},
		{`(test(test_))`, []*token{
			&token{typ: tokenLParen},
			&token{typ: tokenAtom, txt: "test"},
			&token{typ: tokenLParen},
			&token{typ: tokenError}},
		},
		{`(t-e123_s4(test-))`, []*token{
			&token{typ: tokenLParen},
			&token{typ: tokenAtom, txt: "t-e123_s4"},
			&token{typ: tokenLParen},
			&token{typ: tokenError}},
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
	for i, v := range a {
		if v.typ != b[i].typ {
			return false
		} else {
			if v.typ == tokenError {
				return true
			} else if v.txt != v.txt {
				return false
			}
		}
	}
	return true
}
