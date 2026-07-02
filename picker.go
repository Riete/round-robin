package robinx

type ID int64

type Picker[T any] interface {
	Add(item T, weight int64) ID
	Remove(id ID)
	Get(id ID) (*WeightedItem[T], bool)
	Contains(id ID) bool
	All() []*WeightedItem[T]
	Len() int
	Available() int
	SetWeight(id ID, weight int64)
	HasNext() bool
	Next() *WeightedItem[T]
	Clear(f func(*WeightedItem[T]))
	Range(f func(*WeightedItem[T]))
}
