package lisp

import (
	"fmt"
	"strings"
	"testing"
)

type testData struct {
	test     string
	expected []*token
}

func TestComment(t *testing.T) {
	tests := []testData{
		{`; This is a comment
; This too is a comment
;;;;
`, []*token{
			&token{typ: tokenComment, val: "; This is a comment"},
			&token{typ: tokenComment, val: "; This too is a comment"},
			&token{typ: tokenComment, val: ";;;;"},
			&token{typ: tokenEOF},
		}},
	}

	if err := runTokenTest(tests); err != nil {
		t.Error(err)
	}
}

func TestParen(t *testing.T) {
	tests := []testData{
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
			&token{typ: tokenComment, val: ";This is my opening comment"},
			&token{typ: tokenRParen},
			&token{typ: tokenRParen},
			&token{typ: tokenRParen},
			&token{typ: tokenRParen},
			&token{typ: tokenEOF}},
		},
	}
	if err := runTokenTest(tests); err != nil {
		t.Error(err)
	}
}

func TestAtom(t *testing.T) {
	tests := []testData{
		{`(t())`, []*token{
			&token{typ: tokenLParen},
			&token{typ: tokenAtom, val: "t"},
			&token{typ: tokenLParen},
			&token{typ: tokenRParen},
			&token{typ: tokenRParen},
			&token{typ: tokenEOF}},
		},
		{`(test(test foo_roo)) ;This is a comment`, []*token{
			&token{typ: tokenLParen},
			&token{typ: tokenAtom, val: "test"},
			&token{typ: tokenLParen},
			&token{typ: tokenAtom, val: "test"},
			&token{typ: tokenAtom, val: "foo_roo"},
			&token{typ: tokenRParen},
			&token{typ: tokenRParen},
			&token{typ: tokenComment, val: ";This is a comment"},
			&token{typ: tokenEOF}},
		},
		{`(test(test_))`, []*token{
			&token{typ: tokenLParen},
			&token{typ: tokenAtom, val: "test"},
			&token{typ: tokenLParen},
			&token{typ: tokenError}},
		},
		{`(t-e123_s4(test-))`, []*token{
			&token{typ: tokenLParen},
			&token{typ: tokenAtom, val: "t-e123_s4"},
			&token{typ: tokenLParen},
			&token{typ: tokenError}},
		},
		{`test`, []*token{
			&token{typ: tokenAtom, val: "test"},
			&token{typ: tokenEOF}},
		},
	}

	if err := runTokenTest(tests); err != nil {
		t.Error(err)
	}
}

func TestNumber(t *testing.T) {
	tests := []testData{
		{`(123())`, []*token{
			&token{typ: tokenLParen},
			&token{typ: tokenNumber, val: 123},
			&token{typ: tokenLParen},
			&token{typ: tokenRParen},
			&token{typ: tokenRParen},
			&token{typ: tokenEOF}},
		},
		{`(123(4s4
	))`, []*token{
			&token{typ: tokenLParen},
			&token{typ: tokenNumber, val: 123},
			&token{typ: tokenLParen},
			&token{typ: tokenError}},
		},
	}
	if err := runTokenTest(tests); err != nil {
		t.Error(err)
	}
}

func runTokenTest(td []testData) error {
	for _, tst := range td {
		sr := strings.NewReader(tst.test)
		lxr := newLexer(sr)
		var tks []*token
		getTokens(lxr, &tks)
		if cmpTokenSlice(tst.expected, tks) == false {
			return fmt.Errorf("For test string %s\nExpected:\t%v\nGot:\t\t%v\n", tst.test, tst.expected, tks)
		}
	}
	return nil
}

func getTokens(l *lexer, tks *[]*token) {
	t := &token{}
	for t.typ != tokenEOF {
		t = l.next()
		*tks = append(*tks, t)
	}
}

func cmpTokenSlice(a []*token, b []*token) bool {
	for i, v := range a {
		if v.typ != b[i].typ {
			return false
		} else {
			switch v.typ {
			case tokenError:
				//testing done return early
				return true
			case tokenComment, tokenAtom:
				//String val tokens
				x, xok := v.val.(string)
				y, yok := b[i].val.(string)
				if xok != true || yok != true || x != y {
					return false
				}
			case tokenNumber:
				x, xok := v.val.(int)
				y, yok := b[i].val.(int)
				if xok != true || yok != true || x != y {
					return false
				}
			case tokenEOF, tokenLParen, tokenRParen, tokenQuote:
				//No val tokens
				continue
			}
		}
	}
	return true
}
