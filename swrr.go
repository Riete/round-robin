package robinx

import (
	"slices"
	"sync"
)

type SmoothWeightedRR[T any] struct {
	items     []*WeightedItem[T]
	itemIndex map[ID]*WeightedItem[T]
	id        ID
	mu        sync.Mutex
}

func (s *SmoothWeightedRR[T]) Add(item T, weight int64) ID {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.id++
	s.itemIndex[s.id] = &WeightedItem[T]{id: s.id, item: item, weight: weight}
	s.items = append(s.items, s.itemIndex[s.id])
	return s.id
}

func (s *SmoothWeightedRR[T]) Remove(id ID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.itemIndex, id)
	s.items = slices.DeleteFunc(s.items, func(item *WeightedItem[T]) bool {
		return item.id == id
	})
}

func (s *SmoothWeightedRR[T]) Get(id ID) (*WeightedItem[T], bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, exists := s.itemIndex[id]
	return item, exists
}

func (s *SmoothWeightedRR[T]) Contains(id ID) bool {
	_, exists := s.Get(id)
	return exists
}

func (s *SmoothWeightedRR[T]) All() []*WeightedItem[T] {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.items
}

func (s *SmoothWeightedRR[T]) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.items)
}

func (s *SmoothWeightedRR[T]) Available() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	var count int
	for _, item := range s.items {
		if item.weight > 0 {
			count++
		}
	}
	return count
}

func (s *SmoothWeightedRR[T]) SetWeight(id ID, weight int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if item, exists := s.itemIndex[id]; exists {
		item.weight = weight
	}
}

func (s *SmoothWeightedRR[T]) Next() *WeightedItem[T] {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.items) == 0 {
		return nil
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
	if selected.weight == 0 {
		return nil
	}
	selected.current -= totalWeight
	return selected
}

func (s *SmoothWeightedRR[T]) Clear(f func(item *WeightedItem[T])) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if f != nil {
		for _, item := range s.items {
			f(item)
		}
	}
	s.items = make([]*WeightedItem[T], 0)
	clear(s.itemIndex)
}

func (s *SmoothWeightedRR[T]) Range(f func(*WeightedItem[T])) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if f != nil {
		for _, item := range s.items {
			f(item)
		}
	}
}

func NewSmoothWeightedPicker[T any]() Picker[T] {
	return &SmoothWeightedRR[T]{itemIndex: make(map[ID]*WeightedItem[T])}
}
