package rr

import (
	"testing"
)

func TestRoundRobin(t *testing.T) {
	r1 := NewRoundRobinItem("r1")
	r2 := NewRoundRobinItem("r2")
	rr := New(r1, r2)
	r3 := NewRoundRobinItem("r3")
	r4 := NewRoundRobinItem("r4")
	rr.Add(r3, r4)
	for range 10 {
		next := rr.Next()
		t.Log(next.Identity(), next.Data())
	}
	t.Log("====")
	rr1, _ := rr.Get(r1.Identity())
	rr.Remove(rr1.Identity(), r2.Identity())
	for range 10 {
		next := rr.Next()
		t.Log(next.Identity(), next.Data())
	}
}

func TestPop(t *testing.T) {
	r1 := NewRoundRobinItem("r1")
	r2 := NewRoundRobinItem("r2")
	r3 := NewRoundRobinItem("r3")
	r4 := NewRoundRobinItem("r4")
	rr := New(r1, r2, r3, r4)

	x := rr.Pop()
	t.Log(x.Data(), x.Identity())
	x = rr.Pop()
	t.Log(x.Data(), x.Identity())
	for _, i := range rr.All() {
		t.Log(i.Identity(), i.Data())
	}

}
