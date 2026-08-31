package int

import (
	"fmt"
	"testing"
)

// func TestAdder(t *testing.T) {

// 	t.Run("Test 1", func(t *testing.T) {
// 		got := Add(2, 2)
// 		want := 4
// 		if got != want {
// 			t.Errorf("got: %d, want: %d", got, want)
// 		}
// 	})

// 	t.Run("Test 2", func(t *testing.T) {
// 		got := Add(2, 2)
// 		fmt.Println(got)
// 		// Output: 4

// 	})

// }

func ExampleAdd() {
	sum := Add(1, 5)
	fmt.Println(sum)
	// Output: 6
}

func BenchmarkAdd(b *testing.B) {
	for b.Loop() {
		Add(1, 5)
	}
}
