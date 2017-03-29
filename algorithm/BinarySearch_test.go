package algorithm

import (
	"testing"
)

func TestBinarySearch(t *testing.T) {
	low := 20
	high := 100
	for i := low; i <= high; i++ {
		val := i
		n, e := BinarySearch(low, high, func(i int) (int, error) {
			return i - val, nil
		})
		if e != nil {
			t.Fatal(e)
		}
		if n != val {
			t.Fatal(n)
		}
	}

}
