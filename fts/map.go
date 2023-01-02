package fts

import "sync"

type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

type SyncMapper[K comparable, V any] interface {
	Get(key K) (V, bool)
	Put(key K, value V)
	Entries() []Entry[K, V]
}

type HashMap[K comparable, V any] struct {
	lock sync.RWMutex
	mp   map[K]V
}

var _ SyncMapper[int, any] = (*HashMap[int, any])(nil)

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

func (dm *HashMap[K, V]) Put(key K, value V) {
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

func (dm *HashMap[K, V]) Entries() []Entry[K, V] {
	dm.lock.Lock()
	defer dm.lock.Unlock()
	entries := make([]Entry[K, V], len(dm.mp))
	idx := 0
	for k, v := range dm.mp {
		entries[idx] = Entry[K, V]{Key: k, Value: v}
		idx++
	}
	return entries
}

func (dm *HashMap[K, V]) Len() int {
	return len(dm.Entries())
}
