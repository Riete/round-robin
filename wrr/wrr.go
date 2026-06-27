package wrr

import (
	"sync"
)

type WeightedItem[T any] struct {
	item   T
	weight int64
}

func NewWeightedItem[T any](item T, weight int64) *WeightedItem[T] {
	return &WeightedItem[T]{item: item, weight: weight}
}

type WeightedRoundRobin[T any] struct {
	items       []*WeightedItem[T]
	totalWeight int64
	counter     int64
	mu          sync.Mutex
}

// Next 获取下一个任务
func (w *WeightedRoundRobin[T]) Next() T {
	w.mu.Lock()
	defer w.mu.Unlock()

	if len(w.items) == 0 {
		return *new(T)
	}

	w.counter++

	pos := w.counter % w.totalWeight
	if pos == 0 {
		pos = w.totalWeight
	}

	for _, item := range w.items {
		pos -= item.weight
		if pos <= 0 {
			return item.item
		}
	}
	return w.items[len(w.items)-1].item
}

func New[T any](items ...*WeightedItem[T]) *WeightedRoundRobin[T] {
	var totalWeight int64
	for _, item := range items {
		totalWeight += item.weight
	}
	return &WeightedRoundRobin[T]{items: items, totalWeight: totalWeight}
}
