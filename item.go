package robinx

type WeightedItem[T any] struct {
	id          ID
	item        T
	weight      int64
	current     int64
	totalWeight int64
}

func (w *WeightedItem[T]) Identity() ID {
	return w.id
}

func (w *WeightedItem[T]) Item() T {
	return w.item
}

func (w *WeightedItem[T]) CurrentWeight() int64 {
	return w.weight + w.current
}

func (w *WeightedItem[T]) NextWeight() int64 {
	return w.CurrentWeight() - w.totalWeight
}
