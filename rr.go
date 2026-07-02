package robinx

import (
	"slices"
	"sync"
)

type RoundRobin[T any] struct {
	items     []*WeightedItem[T]
	itemIndex map[ID]*WeightedItem[T]
	id        ID
	index     int
	mu        sync.Mutex
}

func (r *RoundRobin[T]) Add(item T, _ int64) ID {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.id++
	r.itemIndex[r.id] = &WeightedItem[T]{id: r.id, item: item}
	r.items = append(r.items, r.itemIndex[r.id])
	return r.id
}

func (r *RoundRobin[T]) Remove(id ID) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.itemIndex, id)
	r.items = slices.DeleteFunc(r.items, func(item *WeightedItem[T]) bool {
		return item.id == id
	})
}

func (r *RoundRobin[T]) Get(id ID) (*WeightedItem[T], bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, exists := r.itemIndex[id]
	return item, exists
}

func (r *RoundRobin[T]) Contains(id ID) bool {
	_, exists := r.Get(id)
	return exists
}

func (r *RoundRobin[T]) All() []*WeightedItem[T] {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.items
}

func (r *RoundRobin[T]) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.items)
}

func (r *RoundRobin[T]) Available() int {
	return r.Len()
}

func (r *RoundRobin[T]) SetWeight(ID, int64) {}

func (r *RoundRobin[T]) HasNext() bool {
	return r.Len() > 0
}

func (r *RoundRobin[T]) Next() *WeightedItem[T] {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.items) == 0 {
		return nil
	}
	if r.index >= len(r.items) {
		r.index = 0
	}
	item := r.items[r.index]
	r.index++
	return item
}

func (r *RoundRobin[T]) Clear(f func(item *WeightedItem[T])) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if f != nil {
		for _, item := range r.items {
			f(item)
		}
	}
	r.items = make([]*WeightedItem[T], 0)
	clear(r.itemIndex)
	r.index = 0
}

func (r *RoundRobin[T]) Range(f func(item *WeightedItem[T])) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if f != nil {
		for _, item := range r.items {
			f(item)
		}
	}
}

func NewRoundRobinPicker[T any]() Picker[T] {
	return &RoundRobin[T]{itemIndex: make(map[ID]*WeightedItem[T])}
}
