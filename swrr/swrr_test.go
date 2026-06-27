package swrr

import (
	"testing"
)

func TestSmoothWeightedRoundRobin(t *testing.T) {
	w1 := NewWeightedItem[string]("a", 4)
	w2 := NewWeightedItem[string]("b", 3)
	w3 := NewWeightedItem[string]("c", 2)
	w4 := NewWeightedItem[string]("d", 1)
	swrr := New[string](w1, w2, w3, w4)
	for range 10 {
		t.Log(swrr.Next())
	}
}
