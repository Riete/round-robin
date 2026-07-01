package swrr

import (
	"slices"
	"sync"
)

type WeightedItem[T any] struct {
	identity    int64
	data        T
	weight      int64
	current     int64
	totalWeight int64
}

func (w *WeightedItem[T]) Identity() int64 {
	return w.identity
}

func (w *WeightedItem[T]) Data() T {
	return w.data
}

func (w *WeightedItem[T]) Weight() int64 {
	return w.weight + w.current
}

func (w *WeightedItem[T]) NextWeight() int64 {
	return w.Weight() - w.totalWeight
}

func NewWeightedItem[T any](data T, weight int64) *WeightedItem[T] {
	return &WeightedItem[T]{data: data, weight: weight}
}

type SmoothWeightedRoundRobin[T any] struct {
	items   []*WeightedItem[T]
	idStart int64
	mu      sync.Mutex
}

func (s *SmoothWeightedRoundRobin[T]) Add(items ...*WeightedItem[T]) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range items {
		s.idStart++
		item.identity = s.idStart
		s.items = append(s.items, item)
	}
}

func (s *SmoothWeightedRoundRobin[T]) Get(identity int64) (*WeightedItem[T], bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range s.items {
		if item.identity == identity {
			return item, true
		}
	}
	return nil, false
}

func (s *SmoothWeightedRoundRobin[T]) Remove(identities ...int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = slices.DeleteFunc(s.items, func(w *WeightedItem[T]) bool {
		return slices.Contains(identities, w.identity)
	})
}

func (s *SmoothWeightedRoundRobin[T]) Next() *WeightedItem[T] {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.items) == 0 {
		return nil
	}
	if len(s.items) == 1 {
		return s.items[0]
	}
	var totalWeight int64
	for _, item := range s.items {
		totalWeight += item.weight
		item.current += item.weight
	}
	selected := s.items[0]
	selected.totalWeight = totalWeight
	for _, item := range s.items[1:] {
		item.totalWeight = totalWeight
		if item.current > selected.current {
			selected = item
		}
	}
	selected.current -= totalWeight
	return selected
}

func (s *SmoothWeightedRoundRobin[T]) Replace(items ...*WeightedItem[T]) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = make([]*WeightedItem[T], 0, len(s.items))
	s.items = append(s.items, items...)
}

func (s *SmoothWeightedRoundRobin[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = make([]*WeightedItem[T], 0, cap(s.items))
}

func (s *SmoothWeightedRoundRobin[T]) All() []*WeightedItem[T] {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.items
}

func (s *SmoothWeightedRoundRobin[T]) SetWeight(identity, weight int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range s.items {
		if item.identity == identity {
			item.weight = weight
			return
		}
	}
}

func (s *SmoothWeightedRoundRobin[T]) SetData(identity int64, data T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range s.items {
		if item.identity == identity {
			item.data = data
			return
		}
	}
}

func (s *SmoothWeightedRoundRobin[T]) SetWeightData(identity, weight int64, data T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range s.items {
		if item.identity == identity {
			item.data = data
			item.weight = weight
			return
		}
	}
}

func New[T any](items ...*WeightedItem[T]) *SmoothWeightedRoundRobin[T] {
	swrr := &SmoothWeightedRoundRobin[T]{}
	swrr.Add(items...)
	return swrr
}
