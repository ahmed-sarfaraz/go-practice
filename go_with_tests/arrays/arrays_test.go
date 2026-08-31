package arrays

import (
	"slices"
	"testing"
)

func TestVariadicFunction(t *testing.T) {
	t.Run("Test Variadic test", func(t *testing.T) {
		got := VariadicAdd([]int{1, 2}, []int{2}, []int{3, 8, 7})
		want := []int{3, 2, 18}
		if !slices.Equal(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	})
}
