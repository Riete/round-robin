package rr

import (
	"sync"
)

type RoundRobin[T any] struct {
	items []T
	index int64
	count int64
	mu    sync.Mutex
}

func (r *RoundRobin[T]) Next() T {
	r.mu.Lock()
	defer r.mu.Unlock()
	item := r.items[r.index%r.count]
	r.index++
	return item
}

func New[T any](items ...T) *RoundRobin[T] {
	return &RoundRobin[T]{items: items, count: int64(len(items))}
}
