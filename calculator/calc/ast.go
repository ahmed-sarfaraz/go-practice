package calc

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrDivideByZero is returned by Eval for a division with a zero divisor.
var ErrDivideByZero = errors.New("division by zero")

// Op is an arithmetic operator. It is deliberately a separate type from the
// lexer's tokenKind: the tree describes meaning, not the characters that
// happened to produce it. The parser translates one into the other, and that
// translation is the only place the two vocabularies meet.
type Op int

const (
	OpAdd Op = iota
	OpSub
	OpMul
	OpDiv
)

func (o Op) String() string {
	switch o {
	case OpAdd:
		return "+"
	case OpSub:
		return "-"
	case OpMul:
		return "*"
	case OpDiv:
		return "/"
	}
	return "?"
}

// Node is one node of the parse tree.
//
// Label and Children exist so that anything walking the tree — the renderer
// below, and whatever you write next — works for every node type without a
// type switch. Add a node type and the compiler will not let you forget to
// implement them, which is exactly the failure a type switch would have hidden.
type Node interface {
	// Eval computes the node's value, recursing into its children.
	Eval() (float64, error)
	// Label is this node's own text, without its children: "+", "neg", "3".
	Label() string
	// Children returns the node's operands in source order, nil for a leaf.
	Children() []Node
}

// Number is a leaf.
type Number struct{ Value float64 }

func (n Number) Eval() (float64, error) { return n.Value, nil }
func (n Number) Children() []Node       { return nil }

// Label prints the value at full precision. A tree dump is a debugging view, so
// it shows what was actually parsed; rounding for human display is the caller's
// job, not the library's.
func (n Number) Label() string { return strconv.FormatFloat(n.Value, 'g', -1, 64) }

// Binary is an operator with two operands: the '+' in 2 + 3.
type Binary struct {
	Op          Op
	Left, Right Node
}

func (b Binary) Label() string    { return b.Op.String() }
func (b Binary) Children() []Node { return []Node{b.Left, b.Right} }

func (b Binary) Eval() (float64, error) {
	left, err := b.Left.Eval()
	if err != nil {
		return 0, err
	}
	right, err := b.Right.Eval()
	if err != nil {
		return 0, err
	}

	switch b.Op {
	case OpAdd:
		return left + right, nil
	case OpSub:
		return left - right, nil
	case OpMul:
		return left * right, nil
	case OpDiv:
		if right == 0 {
			return 0, ErrDivideByZero
		}
		return left / right, nil
	}
	return 0, fmt.Errorf("internal error: unknown operator %d", b.Op)
}

// Negate is the unary minus in -5.
type Negate struct{ Operand Node }

func (n Negate) Label() string    { return "neg" }
func (n Negate) Children() []Node { return []Node{n.Operand} }

func (n Negate) Eval() (float64, error) {
	value, err := n.Operand.Eval()
	if err != nil {
		return 0, err
	}
	return -value, nil
}

// Tree renders a parse tree, so the structure the parser built is visible.
// It is written against the Node interface alone, so it already handles node
// types that do not exist yet — including ones with three or more children.
func Tree(n Node) string {
	var b strings.Builder
	writeTree(&b, n, "", "")
	return b.String()
}

func writeTree(b *strings.Builder, n Node, prefix, childPrefix string) {
	fmt.Fprintf(b, "%s%s\n", prefix, n.Label())

	children := n.Children()
	for i, child := range children {
		branch, indent := "├── ", "│   "
		if i == len(children)-1 {
			branch, indent = "└── ", "    "
		}
		writeTree(b, child, childPrefix+branch, childPrefix+indent)
	}
}
