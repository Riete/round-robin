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

func (r *RoundRobin[T]) Get(identity int64) (*RoundRobinItem[T], bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, item := range r.items {
		if item.identity == identity {
			return item, true
		}
	}
	return nil, false
}

func (r *RoundRobin[T]) Remove(identities ...int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items = slices.DeleteFunc(r.items, func(r *RoundRobinItem[T]) bool {
		return slices.Contains(identities, r.identity)
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

func (r *RoundRobin[T]) Replace(items ...*RoundRobinItem[T]) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items = make([]*RoundRobinItem[T], 0, len(r.items))
	r.items = append(r.items, items...)
}

func (r *RoundRobin[T]) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items = make([]*RoundRobinItem[T], 0, cap(r.items))
}

func (r *RoundRobin[T]) All() []*RoundRobinItem[T] {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.items
}

func (r *RoundRobin[T]) Update(identity int64, data T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, item := range r.items {
		if item.identity == identity {
			item.data = data
			return
		}
	}
}

func New[T any](items ...*RoundRobinItem[T]) *RoundRobin[T] {
	rr := new(RoundRobin[T])
	rr.Add(items...)
	return rr
}
