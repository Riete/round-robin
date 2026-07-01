package robinx

type ID int64

type Picker[T any] interface {
	Add(item T, weight int64) ID
	Remove(id ID)
	Get(id ID) (*WeightedItem[T], bool)
	Contains(id ID) bool
	All() []*WeightedItem[T]
	Len() int
	SetWeight(id ID, weight int64)
	Next() *WeightedItem[T]
	Clear(func(*WeightedItem[T]))
	Range(func(*WeightedItem[T]))
}
