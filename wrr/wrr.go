package wrr

import (
	"slices"
	"sync"
)

type WeightedItem[T any] struct {
	identity int64
	data     T
	weight   int64
}

func (w *WeightedItem[T]) Identity() int64 {
	return w.identity
}

func (w *WeightedItem[T]) Data() T {
	return w.data
}

func (w *WeightedItem[T]) Weight() int64 {
	return w.weight
}

func NewWeightedItem[T any](data T, weight int64) *WeightedItem[T] {
	return &WeightedItem[T]{data: data, weight: weight}
}

type WeightedRoundRobin[T any] struct {
	items       []*WeightedItem[T]
	totalWeight int64
	idStart     int64
	counter     int64
	mu          sync.Mutex
}

func (w *WeightedRoundRobin[T]) Add(items ...*WeightedItem[T]) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, item := range items {
		w.idStart++
		item.identity = w.idStart
		w.items = append(w.items, item)
	}
	for _, item := range items {
		w.totalWeight += item.weight
	}
}

func (w *WeightedRoundRobin[T]) Remove(items ...*WeightedItem[T]) {
	w.mu.Lock()
	defer w.mu.Unlock()
	itemIds := make([]int64, 0, len(items))
	for _, item := range items {
		itemIds = append(itemIds, item.identity)
	}
	w.items = slices.DeleteFunc(w.items, func(w *WeightedItem[T]) bool {
		return slices.Contains(itemIds, w.identity)
	})
	for _, item := range items {
		w.totalWeight -= item.weight
	}
}

func (w *WeightedRoundRobin[T]) Next() *WeightedItem[T] {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.items) == 0 {
		return nil
	}
	if len(w.items) == 1 {
		return w.items[0]
	}
	w.counter++
	pos := w.counter % w.totalWeight
	if pos == 0 {
		pos = w.totalWeight
	}
	for _, item := range w.items {
		pos -= item.weight
		if pos <= 0 {
			return item
		}
	}
	return w.items[len(w.items)-1]
}

func New[T any](items ...*WeightedItem[T]) *WeightedRoundRobin[T] {
	wrr := &WeightedRoundRobin[T]{}
	wrr.Add(items...)
	return wrr
}
