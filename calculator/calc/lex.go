package calc

import (
	"fmt"
	"strconv"
)

// SyntaxError reports a problem with the input text. Pos is a 0-based byte
// offset so a caller can point a caret at the exact column; the message renders
// it 1-based, the way editors count.
type SyntaxError struct {
	Pos int
	Msg string
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("%s (column %d)", e.Msg, e.Pos+1)
}

// tokenKind and token are unexported: they are the lexer's private vocabulary.
// Nothing outside this package — and nothing in ast.go — depends on them.
type tokenKind int

const (
	tokEOF tokenKind = iota
	tokNumber
	tokPlus
	tokMinus
	tokStar
	tokSlash
	tokLParen
	tokRParen
)

// String renders a token the way it should appear inside an error message,
// quoted: "expected ')'". The tree renderer uses Op.String instead, so this
// formatting choice stays local to error reporting.
func (k tokenKind) String() string {
	switch k {
	case tokEOF:
		return "end of input"
	case tokNumber:
		return "a number"
	case tokPlus:
		return "'+'"
	case tokMinus:
		return "'-'"
	case tokStar:
		return "'*'"
	case tokSlash:
		return "'/'"
	case tokLParen:
		return "'('"
	case tokRParen:
		return "')'"
	}
	return "unknown token"
}

type token struct {
	kind tokenKind
	num  float64 // set only when kind == tokNumber
	pos  int
}

var punctuation = map[byte]tokenKind{
	'+': tokPlus, '-': tokMinus, '*': tokStar, '/': tokSlash,
	'(': tokLParen, ')': tokRParen,
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

// lex groups raw characters into tokens: "12 + 3" becomes
// [number(12)] ['+'] [number(3)] [EOF]. This is the step that knows "12" is a
// single number rather than a 1 and a 2.
func lex(input string) ([]token, error) {
	var tokens []token

	for i := 0; i < len(input); {
		c := input[i]

		switch {
		case c == ' ' || c == '\t':
			i++

		case isDigit(c) || c == '.':
			start := i
			for i < len(input) && (isDigit(input[i]) || input[i] == '.') {
				i++
			}
			text := input[start:i]
			value, err := strconv.ParseFloat(text, 64)
			if err != nil {
				return nil, &SyntaxError{Pos: start, Msg: fmt.Sprintf("%q is not a valid number", text)}
			}
			tokens = append(tokens, token{kind: tokNumber, num: value, pos: start})

		default:
			kind, ok := punctuation[c]
			if !ok {
				return nil, &SyntaxError{Pos: i, Msg: fmt.Sprintf("unexpected character %q", string(c))}
			}
			tokens = append(tokens, token{kind: kind, pos: i})
			i++
		}
	}

	// A sentinel end token means the parser never has to bounds-check.
	return append(tokens, token{kind: tokEOF, pos: len(input)}), nil
}
