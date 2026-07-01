package rr

import (
	"slices"
	"sync"
)

type RoundRobinItem[T any] struct {
	identity int64
	data     T
}

func (r *RoundRobinItem[T]) Identity() int64 {
	return r.identity
}

func (r *RoundRobinItem[T]) Data() T {
	return r.data
}

func NewRoundRobinItem[T any](data T) *RoundRobinItem[T] {
	return &RoundRobinItem[T]{data: data}
}

type RoundRobin[T any] struct {
	items   []*RoundRobinItem[T]
	idStart int64
	index   int64
	count   int64
	mu      sync.Mutex
}

func (r *RoundRobin[T]) Add(items ...*RoundRobinItem[T]) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, item := range items {
		r.idStart++
		item.identity = r.idStart
		r.items = append(r.items, item)
	}
	r.count = int64(len(r.items))
}

func (r *RoundRobin[T]) Remove(items ...*RoundRobinItem[T]) {
	r.mu.Lock()
	defer r.mu.Unlock()
	itemIds := make([]int64, 0, len(items))
	for _, item := range items {
		itemIds = append(itemIds, item.identity)
	}
	r.items = slices.DeleteFunc(r.items, func(r *RoundRobinItem[T]) bool {
		return slices.Contains(itemIds, r.identity)
	})
	r.count = int64(len(r.items))
}

func (r *RoundRobin[T]) Next() *RoundRobinItem[T] {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.items) == 0 {
		return nil
	}
	if len(r.items) == 1 {
		return r.items[0]
	}
	if r.index >= r.count {
		r.index = 0
	}
	item := r.items[r.index]
	r.index++
	return item
}

func New[T any](items ...*RoundRobinItem[T]) *RoundRobin[T] {
	rr := new(RoundRobin[T])
	rr.Add(items...)
	return rr
}
