package rr

import (
	"testing"
)

func TestRoundRobin(t *testing.T) {
	rr := New[string]("a", "b")
	for range 10 {
		t.Log(rr.Next())
	}
}
