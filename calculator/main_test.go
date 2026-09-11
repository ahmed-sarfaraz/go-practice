package main

import (
	"strings"
	"testing"

	"calculator/calc"
)

// Display rounding is a main-package concern now, so it is tested here.
func TestFormatNumberHidesFloatArtefacts(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"0.1 + 0.2", "0.3"}, // really 0.30000000000000004 in float64
		{"6 * 7", "42"},      // no trailing ".0"
		{"7 / 2", "3.5"},     // and no integer truncation
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			value, err := calc.Eval(tt.input)
			if err != nil {
				t.Fatalf("Eval(%q) returned unexpected error: %v", tt.input, err)
			}
			if got := formatNumber(value); got != tt.want {
				t.Fatalf("formatNumber(Eval(%q)) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// The caret in an error report has to land under the offending column.
func TestReportPointsAtTheError(t *testing.T) {
	_, err := calc.Parse("2 + * 4")
	if err == nil {
		t.Fatal("Parse succeeded, want a syntax error")
	}

	var out strings.Builder
	report(&out, "2 + * 4", err)

	lines := strings.Split(strings.TrimRight(out.String(), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("report wrote %d lines, want 2:\n%s", len(lines), out.String())
	}
	if caret := strings.Index(lines[1], "^"); caret != strings.Index(lines[0], "*") {
		t.Fatalf("caret at %d, want it under the '*' at %d:\n%s",
			caret, strings.Index(lines[0], "*"), out.String())
	}
}

func TestEvaluateReportsFailure(t *testing.T) {
	var out, errOut strings.Builder

	if evaluate(&out, &errOut, "8 / 0", false) {
		t.Fatal("evaluate returned true for a division by zero")
	}
	if !strings.Contains(errOut.String(), "division by zero") {
		t.Fatalf("stderr = %q, want it to mention division by zero", errOut.String())
	}
	if out.String() != "" {
		t.Fatalf("stdout = %q, want nothing written on failure", out.String())
	}
}
