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

func (w *WeightedRoundRobin[T]) Get(identity int64) (*WeightedItem[T], bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, item := range w.items {
		if item.identity == identity {
			return item, true
		}
	}
	return nil, false
}

func (w *WeightedRoundRobin[T]) Remove(identities ...int64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.items = slices.DeleteFunc(w.items, func(item *WeightedItem[T]) bool {
		if slices.Contains(identities, item.identity) {
			w.totalWeight -= item.weight
			return true
		}
		return false
	})
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

func (w *WeightedRoundRobin[T]) Replace(items ...*WeightedItem[T]) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.items = make([]*WeightedItem[T], 0, len(w.items))
	w.items = append(w.items, items...)
}

func (w *WeightedRoundRobin[T]) Clear() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.items = make([]*WeightedItem[T], 0, cap(w.items))
}

func (w *WeightedRoundRobin[T]) All() []*WeightedItem[T] {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.items
}

func (w *WeightedRoundRobin[T]) SetWeight(identity, weight int64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, item := range w.items {
		if item.identity == identity {
			w.totalWeight += weight - item.weight
			item.weight = weight
			return
		}
	}
}

func (w *WeightedRoundRobin[T]) SetData(identity int64, data T) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, item := range w.items {
		if item.identity == identity {
			item.data = data
			return
		}
	}
}

func (w *WeightedRoundRobin[T]) SetWeightData(identity, weight int64, data T) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, item := range w.items {
		if item.identity == identity {
			w.totalWeight += weight - item.weight
			item.weight = weight
			item.data = data
			return
		}
	}
}

func (w *WeightedRoundRobin[T]) Pop() *WeightedItem[T] {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.items) == 0 {
		return nil
	}
	item := w.items[0]
	w.totalWeight -= item.weight
	w.items = slices.Delete(w.items, 0, 1)
	return item
}

func New[T any](items ...*WeightedItem[T]) *WeightedRoundRobin[T] {
	wrr := &WeightedRoundRobin[T]{}
	wrr.Add(items...)
	return wrr
}
