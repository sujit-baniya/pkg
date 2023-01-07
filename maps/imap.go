package maps

type IMap[K comparable, V any] interface {
	Del(keys ...K)
	Get(key K) (value V, ok bool)
	Set(key K, value V)
	GetOrSet(key K, value V) (actual V, loaded bool)
	GetOrCompute(key K, valueFn func() V) (actual V, loaded bool)
	GetAndDel(key K) (value V, ok bool)
	CompareAndSwap(key K, oldValue, newValue V) bool
	Swap(key K, newValue V) (oldValue V, swapped bool)
	ForEach(lambda func(K, V) bool)
	Grow(newSize uintptr)
	SetHasher(hs func(K) uintptr)
	Len() uintptr
	FillRate() uintptr
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(i []byte) error
}
