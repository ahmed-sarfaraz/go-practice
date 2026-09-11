// Package calc parses and evaluates arithmetic expressions using a recursive
// descent parser. The grammar it recognises is:
//
//	expression → term (('+' | '-') term)*
//	term       → factor (('*' | '/') factor)*
//	factor     → NUMBER | '(' expression ')' | '-' factor
//
// Each rule has a function of the same name below. Precedence is not written
// down anywhere: it comes from term() being called by expression(), so term
// consumes all the '*' and '/' before expression ever sees them.
package calc

import "fmt"

// binaryOps is the boundary between the lexer's vocabulary and the tree's.
// It is the only place a tokenKind becomes an Op.
var binaryOps = map[tokenKind]Op{
	tokPlus:  OpAdd,
	tokMinus: OpSub,
	tokStar:  OpMul,
	tokSlash: OpDiv,
}

type parser struct {
	tokens []token
	pos    int
}

func (p *parser) peek() token { return p.tokens[p.pos] }

// match consumes and returns the next token if it is one of kinds.
func (p *parser) match(kinds ...tokenKind) (token, bool) {
	for _, kind := range kinds {
		if p.peek().kind == kind {
			t := p.tokens[p.pos]
			p.pos++
			return t, true
		}
	}
	return token{}, false
}

// expression → term (('+' | '-') term)*
func (p *parser) expression() (Node, error) {
	node, err := p.term()
	if err != nil {
		return nil, err
	}

	for {
		op, ok := p.match(tokPlus, tokMinus)
		if !ok {
			return node, nil // the (...)* loop ends when no operator matches
		}
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		// Left-associative: 1-2-3 becomes ((1-2)-3), not (1-(2-3)).
		node = Binary{Op: binaryOps[op.kind], Left: node, Right: right}
	}
}

// term → factor (('*' | '/') factor)*
//
// Identical in shape to expression, one precedence level tighter. Because
// expression calls this, '3 * 4' is fully consumed here and handed back as a
// single node before expression gets to look for '+'.
func (p *parser) term() (Node, error) {
	node, err := p.factor()
	if err != nil {
		return nil, err
	}

	for {
		op, ok := p.match(tokStar, tokSlash)
		if !ok {
			return node, nil
		}
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		node = Binary{Op: binaryOps[op.kind], Left: node, Right: right}
	}
}

// factor → NUMBER | '(' expression ')' | '-' factor
func (p *parser) factor() (Node, error) {
	if t, ok := p.match(tokNumber); ok {
		return Number{Value: t.num}, nil
	}

	// '-' factor, so -5 and -(2+3) both work.
	if _, ok := p.match(tokMinus); ok {
		operand, err := p.factor()
		if err != nil {
			return nil, err
		}
		return Negate{Operand: operand}, nil
	}

	// '(' expression ')' — the call back up to expression is what allows
	// nesting to any depth, and what makes this "recursive" descent.
	if _, ok := p.match(tokLParen); ok {
		inner, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, ok := p.match(tokRParen); !ok {
			return nil, &SyntaxError{Pos: p.peek().pos, Msg: "expected ')'"}
		}
		return inner, nil
	}

	found := p.peek()
	return nil, &SyntaxError{
		Pos: found.pos,
		Msg: fmt.Sprintf("expected a number, '(' or '-', found %s", found.kind),
	}
}

// Parse turns an expression into a tree without evaluating it.
func Parse(input string) (Node, error) {
	tokens, err := lex(input)
	if err != nil {
		return nil, err
	}

	p := &parser{tokens: tokens}
	node, err := p.expression()
	if err != nil {
		return nil, err
	}

	// expression() stops at the first token it cannot use. If that is not the
	// end of input, there is trailing junk: "2 3" or "(1+2))".
	if extra := p.peek(); extra.kind != tokEOF {
		return nil, &SyntaxError{Pos: extra.pos, Msg: fmt.Sprintf("unexpected %s", extra.kind)}
	}
	return node, nil
}

// Eval parses and evaluates an expression in one step.
func Eval(input string) (float64, error) {
	node, err := Parse(input)
	if err != nil {
		return 0, err
	}
	return node.Eval()
}
