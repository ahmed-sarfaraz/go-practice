// Command calculator evaluates arithmetic expressions.
//
//	calculator '2 + 3 * 4'      evaluate one expression
//	calculator -tree '2 + 3'    show the parse tree too
//	calculator                  read expressions from stdin, one per line
//
// All parsing and evaluation lives in the calc package; this file is only
// input, output and exit codes.
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"calculator/calc"
)

func main() {
	showTree := flag.Bool("tree", false, "print the parse tree as well as the result")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: calculator [-tree] [expression]\n\n")
		fmt.Fprintf(os.Stderr, "  calculator '2 + 3 * 4'   evaluate one expression\n")
		fmt.Fprintf(os.Stderr, "  calculator               read expressions from stdin, one per line\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	// An expression on the command line: evaluate it and exit. This is what
	// makes the tool usable from a script or a pipeline.
	if flag.NArg() > 0 {
		if !evaluate(os.Stdout, os.Stderr, strings.Join(flag.Args(), " "), *showTree) {
			os.Exit(1)
		}
		return
	}

	if !repl(os.Stdin, os.Stdout, os.Stderr, *showTree) {
		os.Exit(1)
	}
}

// formatNumber trims to 12 significant digits for display. float64 cannot
// represent 0.1 exactly, so 0.1+0.2 is really 0.30000000000000004; rounding
// here hides that artefact without touching the arithmetic.
//
// This is a presentation decision, which is why it lives in main and not in
// calc: the library returns a float64 and lets each caller render it.
func formatNumber(v float64) string {
	return strconv.FormatFloat(v, 'g', 12, 64)
}

// repl reads one expression per line. It returns false if any line failed and
// the input was not a terminal, so piped input still sets a useful exit code.
func repl(in *os.File, out, errOut io.Writer, showTree bool) bool {
	interactive := isTerminal(in)
	if interactive {
		fmt.Fprintln(out, "Type an expression, or Ctrl-D to quit.")
	}

	scanner := bufio.NewScanner(in)
	allOK := true

	for {
		if interactive {
			fmt.Fprint(out, "> ")
		}
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if !evaluate(out, errOut, line, showTree) {
			allOK = false
		}
	}

	if interactive {
		fmt.Fprintln(out)
		return true // a typo during a session is not a failed run
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(errOut, "calculator: %v\n", err)
		return false
	}
	return allOK
}

func evaluate(out, errOut io.Writer, line string, showTree bool) bool {
	node, err := calc.Parse(line)
	if err != nil {
		report(errOut, line, err)
		return false
	}

	if showTree {
		fmt.Fprint(out, calc.Tree(node))
	}

	value, err := node.Eval()
	if err != nil {
		fmt.Fprintf(errOut, "calculator: %v\n", err)
		return false
	}

	fmt.Fprintln(out, formatNumber(value))
	return true
}

// report points a caret at the offending column, which is the whole reason
// calc.SyntaxError carries a position.
func report(w io.Writer, line string, err error) {
	var syntaxErr *calc.SyntaxError
	if errors.As(err, &syntaxErr) {
		pad := min(syntaxErr.Pos, len(line))
		fmt.Fprintf(w, "  %s\n  %s^ %s\n", line, strings.Repeat(" ", pad), syntaxErr.Msg)
		return
	}
	fmt.Fprintf(w, "calculator: %v\n", err)
}

func isTerminal(f *os.File) bool {
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
