package robinx

type WeightedItem[T any] struct {
	id          ID
	item        T
	weight      int64
	current     int64
	totalWeight int64
}

func (w *WeightedItem[T]) ID() ID {
	return w.id
}

func (w *WeightedItem[T]) Item() T {
	return w.item
}

func (w *WeightedItem[T]) SetWeight(weight int64) {
	w.weight = weight
}

func (w *WeightedItem[T]) Weight() int64 {
	return w.weight
}

func (w *WeightedItem[T]) CurrentWeight() int64 {
	return w.weight + w.current
}

func (w *WeightedItem[T]) NextWeight() int64 {
	return w.CurrentWeight() - w.totalWeight
}

func (w *WeightedItem[T]) TotalWeight() int64 {
	return w.totalWeight
}
