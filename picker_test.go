package robinx

import (
	"fmt"
	"testing"
)

var rr = NewRoundRobinPicker[string]()
var swrr = NewSmoothWeightedPicker[string]()

var f = func(item *WeightedItem[string]) {
	if item == nil {
		fmt.Println("nil")
	} else {
		fmt.Println("id", item.id, "weight", item.weight, "data", item.item)
	}
}

func TestPickerAdd(t *testing.T) {
	rr.Add("a", 0)
	rr.Add("b", 0)
	rr.Range(f)

	swrr.Add("a", 1)
	swrr.Add("b", 1)
	swrr.Range(f)
}

func TestPickerRemove(t *testing.T) {
	r1 := rr.Add("a", 0)
	rr.Add("b", 0)
	rr.Remove(r1)
	rr.Range(f)

	s1 := swrr.Add("a", 1)
	swrr.Add("b", 1)
	swrr.Remove(s1)
	swrr.Range(f)
}

func TestPickerGet(t *testing.T) {
	r1 := rr.Add("a", 0)
	t.Log(rr.Get(r1))

	s1 := swrr.Add("a", 1)
	t.Log(swrr.Get(s1))
}

func TestPickerContains(t *testing.T) {
	r1 := rr.Add("a", 0)
	t.Log(rr.Contains(r1))

	s1 := swrr.Add("a", 1)
	t.Log(swrr.Contains(s1))
}

func TestPickerAll(t *testing.T) {
	rr.Add("a", 0)
	rr.Add("b", 0)
	for _, item := range rr.All() {
		f(item)
	}

	swrr.Add("a", 1)
	swrr.Add("b", 1)
	for _, item := range swrr.All() {
		f(item)
	}
}

func TestPickerSetWeight(t *testing.T) {
	s1 := swrr.Add("a", 1)
	s2 := swrr.Add("b", 1)
	swrr.SetWeight(s1, 100)
	swrr.SetWeight(s2, 12)
	for _, item := range swrr.All() {
		f(item)
	}
}

func TestPickerNext(t *testing.T) {
	s1 := swrr.Add("a", 9)
	s2 := swrr.Add("b", 1)
	for range 10 {
		f(swrr.Next())
	}
	swrr.SetWeight(s1, 0)
	swrr.SetWeight(s2, 1)
	t.Log("===")
	for range 10 {
		f(swrr.Next())
	}
}

func TestPickerClear(t *testing.T) {
	rr.Add("a", 0)
	rr.Add("b", 0)
	rr.Clear(f)
	t.Log(rr.All())

	swrr.Add("a", 1)
	swrr.Add("b", 1)
	swrr.Clear(f)
	t.Log(swrr.All())
}
