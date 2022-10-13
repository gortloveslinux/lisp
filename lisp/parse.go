package lisp

type parseError error

type parser struct {
	l *lexer
}

func newParser(l *lexer) *parser {
	return &parser{l: l}
}

type expr struct {
	first *expr
	atom  *token
	rest  *expr
}

func (p *parser) parseSExpr() (*expr, error) {
	for {
		t := p.l.next()
		switch t.typ {
		case tokenLParen:
			//f, err := p.parseSExpr()
			t = p.l.next()
			if t.typ != tokenRParen {
				//return &expr{first: f	}, fmt.Errorf("Expecting ')' encountered '%v'", t.val
			}
		}
	}
}
