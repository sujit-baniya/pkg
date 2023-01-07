package fts

import (
	"github.com/goccy/go-reflect"
	"sync"
)

type HashMap[K comparable, V any] struct {
	lock sync.RWMutex
	mp   map[K]V
}

func NewMap[K comparable, V any]() *HashMap[K, V] {
	return &HashMap[K, V]{
		mp: make(map[K]V),
	}
}

func (dm *HashMap[K, V]) Get(key K) (V, bool) {
	dm.lock.Lock()
	defer dm.lock.Unlock()
	val, ok := dm.mp[key]
	return val, ok
}

func (dm *HashMap[K, V]) Set(key K, value V) {
	dm.lock.Lock()
	defer dm.lock.Unlock()
	dm.mp[key] = value
}

func (dm *HashMap[K, V]) GetOrInsert(key K, value V) (V, bool) {
	dm.lock.Lock()
	defer dm.lock.Unlock()
	val, ok := dm.mp[key]
	if ok {
		return val, ok
	}
	dm.mp[key] = value
	return value, ok
}

func (dm *HashMap[K, V]) Del(key K) {
	dm.lock.Lock()
	defer dm.lock.Unlock()
	delete(dm.mp, key)
}

func (dm *HashMap[K, V]) Entries() []map[K]V {
	dm.lock.Lock()
	defer dm.lock.Unlock()
	entries := make([]map[K]V, len(dm.mp))
	idx := 0
	for k, v := range dm.mp {
		entries[idx] = map[K]V{k: v}
		idx++
	}
	return entries
}

func (dm *HashMap[K, V]) Len() int {
	return len(dm.Entries())
}

// Keys returns a slice of the map's keys
func (dm *HashMap[K, V]) Keys() []K {
	keys := make([]K, len(dm.mp))

	var i int
	for k := range dm.mp {
		keys[i] = k
		i++
	}

	return keys
}

// Values returns a slice of the map's values
func (dm *HashMap[K, V]) Values() []V {
	values := make([]V, len(dm.mp))

	var i int
	for _, v := range dm.mp {
		values[i] = v
		i++
	}

	return values
}

// Merge returns a slice of the map's values
func (dm *HashMap[K, V]) Merge(maps ...map[K]V) map[K]V {
	result := make(map[K]V, len(dm.mp)+len(maps))
	for k, v := range dm.mp {
		result[k] = v
	}
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	dm.mp = result
	return result
}

// ForEach returns a slice of the map's values
func (dm *HashMap[K, V]) ForEach(iteratee func(key K, value V)) {
	for k, v := range dm.mp {
		iteratee(k, v)
	}
}

// Filter returns a slice of the map's values
func (dm *HashMap[K, V]) Filter(predicate func(key K, value V) bool) map[K]V {
	result := make(map[K]V)

	for k, v := range dm.mp {
		if predicate(k, v) {
			result[k] = v
		}
	}
	return result
}

// Intersect returns a slice of the map's values
func (dm *HashMap[K, V]) Intersect(maps ...map[K]V) map[K]V {
	if len(maps) == 0 {
		return dm.mp
	}

	var result map[K]V

	reducer := func(m1, m2 map[K]V) map[K]V {
		m := make(map[K]V)
		for k, v1 := range m1 {
			if v2, ok := m2[k]; ok && reflect.DeepEqual(v1, v2) {
				m[k] = v1
			}
		}
		return m
	}

	reduceMaps := make([]map[K]V, 2)
	result = reducer(dm.mp, maps[0])

	for i := 1; i < len(maps); i++ {
		reduceMaps[0] = result
		reduceMaps[1] = maps[i]
		result = reducer(reduceMaps[0], reduceMaps[1])
	}

	return result
}

// Minus returns a slice of the map's values
func (dm *HashMap[K, V]) Minus(mapB map[K]V) map[K]V {
	result := make(map[K]V)

	for k, v := range dm.mp {
		if _, ok := mapB[k]; !ok {
			result[k] = v
		}
	}
	return result
}

// IsDisjoint returns a slice of the map's values
func (dm *HashMap[K, V]) IsDisjoint(mapB map[K]V) bool {
	for k := range dm.mp {
		if _, ok := mapB[k]; ok {
			return false
		}
	}
	return true
}

// Merge maps, next key will overwrite previous key
func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
	result := make(map[K]V, 0)

	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}

	return result
}

// ForEach executes iteratee function for every key and value pair in map
func ForEach[K comparable, V any](m map[K]V, iteratee func(key K, value V)) {
	for k, v := range m {
		iteratee(k, v)
	}
}

// Filter iterates over map, return a new map contains all key and value pairs pass the predicate function
func Filter[K comparable, V any](m map[K]V, predicate func(key K, value V) bool) map[K]V {
	result := make(map[K]V)

	for k, v := range m {
		if predicate(k, v) {
			result[k] = v
		}
	}
	return result
}

// Intersect iterates over maps, return a new map of key and value pairs in all given maps
func Intersect[K comparable, V any](maps ...map[K]V) map[K]V {
	if len(maps) == 0 {
		return map[K]V{}
	}
	if len(maps) == 1 {
		return maps[0]
	}

	var result map[K]V

	reducer := func(m1, m2 map[K]V) map[K]V {
		m := make(map[K]V)
		for k, v1 := range m1 {
			if v2, ok := m2[k]; ok && reflect.DeepEqual(v1, v2) {
				m[k] = v1
			}
		}
		return m
	}

	reduceMaps := make([]map[K]V, 2)
	result = reducer(maps[0], maps[1])

	for i := 2; i < len(maps); i++ {
		reduceMaps[0] = result
		reduceMaps[1] = maps[i]
		result = reducer(reduceMaps[0], reduceMaps[1])
	}

	return result
}

// Minus creates an map of whose key in mapA but not in mapB
func Minus[K comparable, V any](mapA, mapB map[K]V) map[K]V {
	result := make(map[K]V)

	for k, v := range mapA {
		if _, ok := mapB[k]; !ok {
			result[k] = v
		}
	}
	return result
}

// IsDisjoint two map are disjoint if they have no keys in common
func IsDisjoint[K comparable, V any](mapA, mapB map[K]V) bool {
	for k := range mapA {
		if _, ok := mapB[k]; ok {
			return false
		}
	}
	return true
}
