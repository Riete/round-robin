package swrr

import (
	"sync"
)

type WeightedItem[T any] struct {
	item    T
	weight  int64
	current int64
}

func NewWeightedItem[T any](item T, weight int64) *WeightedItem[T] {
	return &WeightedItem[T]{item: item, weight: weight}
}

type SmoothWeightedRoundRobin[T any] struct {
	items []*WeightedItem[T]
	mu    sync.Mutex
}

func (w *SmoothWeightedRoundRobin[T]) Next() T {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.items) == 0 {
		return *new(T)
	}

	var totalWeight int64
	for _, item := range w.items {
		totalWeight += item.weight
		item.current += item.weight
	}

	selected := w.items[0]
	for _, item := range w.items[1:] {
		if item.current > selected.current {
			selected = item
		}
	}
	selected.current -= totalWeight
	return selected.item
}

func New[T any](items ...*WeightedItem[T]) *SmoothWeightedRoundRobin[T] {
	return &SmoothWeightedRoundRobin[T]{items: items}
}
