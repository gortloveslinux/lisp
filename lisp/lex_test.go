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
			&token{typ: tokenComment, val: nil, raw: "; This is a comment", row: 1, col: 1},
			&token{typ: tokenComment, val: nil, raw: "; This too is a comment", row: 2, col: 1},
			&token{typ: tokenComment, val: nil, raw: ";;;;", row: 3, col: 1},
			&token{typ: tokenEOF, row: 4, col: 1},
		}},
	}

	if err := runTokenTest(tests); err != nil {
		t.Error(err)
	}
}

func TestRowCol(t *testing.T) {
	tests := []testData{
		{`; This is a comment
; This too is a comment
;;;;
`, []*token{
			&token{typ: tokenComment, val: nil, raw: "; This is a comment", row: 1, col: 1},
			&token{typ: tokenComment, val: nil, raw: "; This too is a comment", row: 2, col: 1},
			&token{typ: tokenComment, val: nil, raw: ";;;;", row: 3, col: 1},
			&token{typ: tokenEOF, row: 4, col: 1},
		}},
	}

	if err := runTokenTest(tests); err != nil {
		t.Error(err)
	}
}

func TestParen(t *testing.T) {
	tests := []testData{
		{`((()))`, []*token{
			&token{typ: tokenLParen, row: 1, col: 1},
			&token{typ: tokenLParen, row: 1, col: 2},
			&token{typ: tokenLParen, row: 1, col: 3},
			&token{typ: tokenRParen, row: 1, col: 4},
			&token{typ: tokenRParen, row: 1, col: 5},
			&token{typ: tokenRParen, row: 1, col: 6},
			&token{typ: tokenEOF, row: 1, col: 7}},
		},
		{`(;This is my opening comment
))))`, []*token{
			&token{typ: tokenLParen, row: 1, col: 1},
			&token{typ: tokenComment, val: nil, raw: ";This is my opening comment", row: 1, col: 2},
			&token{typ: tokenRParen, row: 2, col: 1},
			&token{typ: tokenRParen, row: 2, col: 2},
			&token{typ: tokenRParen, row: 2, col: 3},
			&token{typ: tokenRParen, row: 2, col: 4},
			&token{typ: tokenEOF, row: 2, col: 5}},
		},
		{`(;This is my opening comment
				))))`, []*token{
			&token{typ: tokenLParen, row: 1, col: 1},
			&token{typ: tokenComment, val: nil, raw: ";This is my opening comment", row: 1, col: 2},
			&token{typ: tokenRParen, row: 2, col: 5},
			&token{typ: tokenRParen, row: 2, col: 6},
			&token{typ: tokenRParen, row: 2, col: 7},
			&token{typ: tokenRParen, row: 2, col: 8},
			&token{typ: tokenEOF, row: 2, col: 9}},
		},
	}
	if err := runTokenTest(tests); err != nil {
		t.Error(err)
	}
}

func TestAtom(t *testing.T) {
	tests := []testData{
		{`(t())`, []*token{
			&token{typ: tokenLParen, row: 1, col: 1},
			&token{typ: tokenAtom, val: "t", raw: "t", row: 1, col: 2},
			&token{typ: tokenLParen, row: 1, col: 3},
			&token{typ: tokenRParen, row: 1, col: 4},
			&token{typ: tokenRParen, row: 1, col: 5},
			&token{typ: tokenEOF, row: 1, col: 6}},
		},
		{`(test(test foo_roo)) ;This is a comment`, []*token{
			&token{typ: tokenLParen, row: 1, col: 1},
			&token{typ: tokenAtom, val: "test", raw: "test", row: 1, col: 2},
			&token{typ: tokenLParen, row: 1, col: 6},
			&token{typ: tokenAtom, val: "test", raw: "test", row: 1, col: 7},
			&token{typ: tokenAtom, val: "foo_roo", raw: "foo_roo", row: 1, col: 12},
			&token{typ: tokenRParen, row: 1, col: 19},
			&token{typ: tokenRParen, row: 1, col: 20},
			&token{typ: tokenComment, val: nil, raw: ";This is a comment", row: 1, col: 22},
			&token{typ: tokenEOF, row: 1, col: 40}},
		},
		{`(test(test_))`, []*token{
			&token{typ: tokenLParen, row: 1, col: 1},
			&token{typ: tokenAtom, val: "test", raw: "test", row: 1, col: 2},
			&token{typ: tokenLParen, row: 1, col: 6},
			&token{typ: tokenError, raw: "test_", row: 1, col: 7}},
		},
		{`(t-e123_s4(test-))`, []*token{
			&token{typ: tokenLParen, row: 1, col: 1},
			&token{typ: tokenAtom, val: "t-e123_s4", raw: "t-e123_s4", row: 1, col: 2},
			&token{typ: tokenLParen, row: 1, col: 11},
			&token{typ: tokenError, row: 1, col: 12, raw: "test-"}},
		},
		{`test`, []*token{
			&token{typ: tokenAtom, val: "test", raw: "test", row: 1, col: 1},
			&token{typ: tokenEOF, row: 1, col: 5}},
		},
	}

	if err := runTokenTest(tests); err != nil {
		t.Error(err)
	}
}

func TestNumber(t *testing.T) {
	tests := []testData{
		{`(123())`, []*token{
			&token{typ: tokenLParen, row: 1, col: 1},
			&token{typ: tokenNumber, val: 123, row: 1, col: 2},
			&token{typ: tokenLParen, row: 1, col: 5},
			&token{typ: tokenRParen, row: 1, col: 6},
			&token{typ: tokenRParen, row: 1, col: 7},
			&token{typ: tokenEOF, row: 1, col: 8}},
		},
		{`(-123())`, []*token{
			&token{typ: tokenLParen, row: 1, col: 1},
			&token{typ: tokenNumber, val: -123, row: 1, col: 2},
			&token{typ: tokenLParen, row: 1, col: 6},
			&token{typ: tokenRParen, row: 1, col: 7},
			&token{typ: tokenRParen, row: 1, col: 8},
			&token{typ: tokenEOF, row: 1, col: 9}},
		},
		{`(-
		123())`, []*token{
			&token{typ: tokenLParen, row: 1, col: 1},
			&token{typ: tokenError, raw: "-", row: 1, col: 2}},
		},
		{`(123(4s4
	))`, []*token{
			&token{typ: tokenLParen, row: 1, col: 1},
			&token{typ: tokenNumber, val: 123, row: 1, col: 2},
			&token{typ: tokenLParen, row: 1, col: 5},
			&token{typ: tokenError, raw: "4s4", row: 1, col: 6}},
		},
	}
	if err := runTokenTest(tests); err != nil {
		t.Error(err)
	}
}

func TestFloat(t *testing.T) {
	tests := []testData{
		{`(123.0)`, []*token{
			&token{typ: tokenLParen, row: 1, col: 1},
			&token{typ: tokenNumber, val: 123.0, row: 1, col: 2},
			&token{typ: tokenRParen, row: 1, col: 7},
			&token{typ: tokenEOF, row: 1, col: 8}},
		},
		{`(-123.0)`, []*token{
			&token{typ: tokenLParen, row: 1, col: 1},
			&token{typ: tokenNumber, val: -123.0, row: 1, col: 2},
			&token{typ: tokenRParen, row: 1, col: 8},
			&token{typ: tokenEOF, row: 1, col: 9}},
		},
		{`(.05)`, []*token{
			&token{typ: tokenLParen, row: 1, col: 1},
			&token{typ: tokenNumber, val: .05, row: 1, col: 2},
			&token{typ: tokenRParen, row: 1, col: 5},
			&token{typ: tokenEOF, row: 1, col: 6}},
		},
		{`(-.05)`, []*token{
			&token{typ: tokenLParen, row: 1, col: 1},
			&token{typ: tokenNumber, val: -.05, row: 1, col: 2},
			&token{typ: tokenRParen, row: 1, col: 6},
			&token{typ: tokenEOF, row: 1, col: 7}},
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
			return fmt.Errorf("For test string %s\nExpected:\t%v\nGot:\t\t\t\t%v\n", tst.test, tst.expected, tks)
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
			case tokenComment, tokenAtom, tokenError:
				x, y := v, b[i]
				if x.raw != y.raw || x.row != y.row || x.col != y.col {
					return false
				}
			case tokenNumber:
				fx, fxok := v.val.(float64)
				fy, fyok := b[i].val.(float64)
				ix, ixok := v.val.(int)
				iy, iyok := b[i].val.(int)
				if fxok && fyok {
					if fx != fy {
						return false
					}
				} else if ixok && iyok {
					if ix != iy {
						return false
					}
				} else {
					return false
				}
			case tokenEOF, tokenLParen, tokenRParen, tokenQuote:
				x, y := v, b[i]
				if x.row != y.row || x.col != y.col {
					return false
				}
				//No val tokens
				continue
			}
		}
	}
	return true
}
