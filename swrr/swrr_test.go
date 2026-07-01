package swrr

import (
	"testing"
)

func TestSmoothWeightedRoundRobin(t *testing.T) {
	s1 := NewWeightedItem[string]("s1", 10)
	s2 := NewWeightedItem[string]("s2", 2)
	s3 := NewWeightedItem[string]("s3", 1)
	s4 := NewWeightedItem[string]("s4", 1)
	s5 := NewWeightedItem[string]("s5", 1)
	swrr := New[string](s1, s2)
	swrr.Add(s3, s4, s5)
	for range 15 {
		next := swrr.Next()
		t.Log(next.Identity(), next.Data(), s1.NextWeight(), s2.NextWeight(), s3.NextWeight(), s4.NextWeight(), s5.NextWeight())
	}
	ss1, _ := swrr.Get(s1.Identity())
	t.Log("====")
	swrr.SetWeight(s5.Identity(), 100)
	swrr.Remove(ss1.Identity(), s2.Identity())
	for range 15 {
		next := swrr.Next()
		t.Log(next.Identity(), next.Data(), s3.NextWeight(), s4.NextWeight(), s5.NextWeight())
	}
}
