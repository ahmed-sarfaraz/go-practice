// Package calc_test is an external test package: it can only reach calc's
// exported API, exactly as a real consumer would. If something here fails to
// compile, the package boundary is wrong.
package calc_test

import (
	"errors"
	"strings"
	"testing"

	"calculator/calc"
)

func TestEval(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		// Precedence: the whole point of the layered grammar.
		{"2 + 3 * 4", 14},
		{"2 * 3 + 4", 10},
		{"10 - 2 * 3", 4},

		// Parentheses override it by pushing the sum down to factor level.
		{"(2 + 3) * 4", 20},
		{"2 * (3 + 4) - 1", 13},
		{"((((5))))", 5},

		// Left associativity: 10-3-2 is (10-3)-2 = 5, not 10-(3-2) = 9.
		{"10 - 3 - 2", 5},
		{"100 / 5 / 2", 10},

		// Unary minus, including stacked and applied to a group.
		{"-5", -5},
		{"-5 + 3", -2},
		{"3 * -2", -6},
		{"--7", 7},
		{"-(2 + 3)", -5},

		// Whitespace is noise; a bare number is a valid expression.
		{"2+3*4", 14},
		{"   7   ", 7},
		{"3.5 * 2", 7},

		// Deeper nesting, mixing everything.
		{"2 + 3 * (4 - 1)", 11},
		{"(1 + 2) * (3 + 4)", 21},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := calc.Eval(tt.input)
			if err != nil {
				t.Fatalf("Eval(%q) returned unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("Eval(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestEvalDivideByZero(t *testing.T) {
	for _, input := range []string{"1 / 0", "5 / (3 - 3)", "2 + 8 / 0"} {
		if _, err := calc.Eval(input); !errors.Is(err, calc.ErrDivideByZero) {
			t.Fatalf("Eval(%q) error = %v, want %v", input, err, calc.ErrDivideByZero)
		}
	}
}

func TestSyntaxErrors(t *testing.T) {
	tests := []struct {
		input   string
		wantPos int    // 0-based column the caret should point at
		wantMsg string // substring of the message
	}{
		{"2 +", 3, "found end of input"},
		{"2 + * 4", 4, "found '*'"},
		{"(2 + 3", 6, "expected ')'"},
		{"2 3", 2, "unexpected a number"},
		{"(1 + 2))", 7, "unexpected ')'"},
		{"2 $ 3", 2, "unexpected character"},
		{"", 0, "found end of input"},
		{"1.2.3 + 1", 0, "not a valid number"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			_, err := calc.Eval(tt.input)
			if err == nil {
				t.Fatalf("Eval(%q) succeeded, want a syntax error", tt.input)
			}

			var syntaxErr *calc.SyntaxError
			if !errors.As(err, &syntaxErr) {
				t.Fatalf("Eval(%q) error = %T, want *calc.SyntaxError", tt.input, err)
			}
			if syntaxErr.Pos != tt.wantPos {
				t.Errorf("Eval(%q) error position = %d, want %d", tt.input, syntaxErr.Pos, tt.wantPos)
			}
			if !strings.Contains(syntaxErr.Msg, tt.wantMsg) {
				t.Errorf("Eval(%q) message = %q, want it to contain %q", tt.input, syntaxErr.Msg, tt.wantMsg)
			}
		})
	}
}

// The tree is the parser's actual output; asserting its shape proves the
// precedence claim rather than just the arithmetic that follows from it.
func TestTreeShape(t *testing.T) {
	node, err := calc.Parse("2 + 3 * 4")
	if err != nil {
		t.Fatalf("Parse returned unexpected error: %v", err)
	}

	want := strings.Join([]string{
		"+",
		"├── 2",
		"└── *",
		"    ├── 3",
		"    └── 4",
		"",
	}, "\n")

	if got := calc.Tree(node); got != want {
		t.Fatalf("Tree =\n%s\nwant\n%s", got, want)
	}
}

// Tree walks the Node interface rather than switching on concrete types, so it
// renders a node type it has never seen — including one with three children.
// This is the property that stops a new node type from silently vanishing.
type triple struct{ a, b, c calc.Node }

func (t triple) Label() string         { return "triple" }
func (t triple) Children() []calc.Node { return []calc.Node{t.a, t.b, t.c} }
func (t triple) Eval() (float64, error) {
	return 0, errors.New("not evaluable")
}

func TestTreeHandlesUnknownNodeTypes(t *testing.T) {
	n := triple{calc.Number{Value: 1}, calc.Number{Value: 2}, calc.Number{Value: 3}}

	want := strings.Join([]string{
		"triple",
		"├── 1",
		"├── 2",
		"└── 3",
		"",
	}, "\n")

	if got := calc.Tree(n); got != want {
		t.Fatalf("Tree =\n%s\nwant\n%s", got, want)
	}
}

// The tree is built from calc.Op, not from lexer tokens, so a caller can
// construct one directly without knowing anything about lexing.
func TestTreeCanBeBuiltByHand(t *testing.T) {
	node := calc.Binary{
		Op:   calc.OpMul,
		Left: calc.Number{Value: 6},
		Right: calc.Binary{
			Op:    calc.OpAdd,
			Left:  calc.Number{Value: 3},
			Right: calc.Number{Value: 4},
		},
	}

	got, err := node.Eval()
	if err != nil {
		t.Fatalf("Eval returned unexpected error: %v", err)
	}
	if got != 42 {
		t.Fatalf("Eval = %v, want 42", got)
	}
}
