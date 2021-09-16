package lisp

import (
	"strings"
	"testing"
)

func TestComment(t *testing.T) {
	l := `; This is a comment
; This too is a comment`

	expected := []*token{
		&token{typ: tokenComment, txt: " This is a comment"},
		&token{typ: tokenComment, txt: " This too is a comment"},
		&token{typ: tokenEOF},
	}

	sr := strings.NewReader(l)
	lxr := newLexer(sr)

	for _, v := range expected {
		tx := lxr.next()
		if !(tx.typ == v.typ && tx.txt == v.txt) {
			t.Errorf("Expected %s to match %s", tx, v)
		}
	}
}
