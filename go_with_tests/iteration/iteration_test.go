package iteration

import (
	"testing"
)

func TestIteration(t *testing.T) {
	t.Run("Concatenate without string package", func(t *testing.T) {
		got := Iterate("s")
		want := "ssss"
		if got != want {
			t.Errorf("got: %s, want %s", got, want)
		}
	})
	t.Run("Concatenate with string package", func(t *testing.T) {
		got := IterateWithStringPackage("s")
		want := "ssss"
		if got != want {
			t.Errorf("got: %s, want %s", got, want)
		}
	})
}

func BenchmarkIteration(b *testing.B) {
	for b.Loop() {
		Iterate("s")
	}
}

func BenchmarkIterationStringPackage(b *testing.B) {
	for b.Loop() {
		IterateWithStringPackage("s")
	}
}
