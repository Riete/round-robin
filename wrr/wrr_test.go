package wrr

import (
	"testing"
)

func TestWeightedRoundRobin(t *testing.T) {
	w1 := NewWeightedItem[string]("w1", 5)
	w2 := NewWeightedItem[string]("w2", 4)
	w3 := NewWeightedItem[string]("w3", 3)
	w4 := NewWeightedItem[string]("w4", 2)
	w5 := NewWeightedItem[string]("w5", 1)

	wrr := New[string](w1, w2, w3, w4)
	wrr.Add(w5)
	for range 15 {
		next := wrr.Next()
		t.Log(next.Identity(), next.Data())
	}
	t.Log("====")
	ww1, _ := wrr.Get(w1.Identity())
	wrr.SetWeight(w3.Identity(), 1)
	wrr.SetWeight(w5.Identity(), 3)
	wrr.Remove(ww1.Identity(), w2.Identity())
	for range 15 {
		next := wrr.Next()
		t.Log(next.Identity(), next.Data())
	}
}
