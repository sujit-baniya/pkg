package fts

import (
	"sort"
)

// Simple has complexity: O(n^2)
func Simple[T comparable](a []T, b []T) []T {
	set := make([]T, 0)

	for _, v := range a {
		if containsGeneric(b, v) {
			set = append(set, v)
		}
	}

	return set
}

// SortedGeneric has complexity: O(n * log(n)), a needs to be sorted
func SortedGeneric[T comparable](a []T, b []T) []T {
	set := make([]T, 0)

	for _, v := range a {
		idx := sort.Search(len(b), func(i int) bool {
			return b[i] == v
		})
		if idx < len(b) && b[idx] == v {
			set = append(set, v)
		}
	}

	return set
}

type Comparator func(i, j int) bool

// Sorted has complexity: O(n + x) where n is length of the shortest array and x duplicate cases in the longest array.
// Best case complexity: O(n) where n is length of the shortest array (all values unique)
// Worst case complexity: O(n) where n is length of the longest array (all values of the longest array are duplicates of intersect match)
// Warning: Function will change left array order
func Sorted[T comparable](a []T, b []T, leftGreater Comparator) []T {
	var i, j, k int

	for {
		if i >= len(a) || j >= len(b) {
			break
		}
		if a[i] == b[j] {
			a[k], a[i] = a[i], a[k]
			i++
			j++
			k++
			continue
		}
		if leftGreater(i, j) {
			j++
			continue
		}
		i++
		continue
	}
	return a[:k]
}

// HashGeneric has complexity: O(n * x) where x is a factor of hash function efficiency (between 1 and 2)
func HashGeneric[T comparable](a []T, b []T) []T {
	set := make([]T, 0)
	hash := make(map[T]struct{})

	for _, v := range a {
		hash[v] = struct{}{}
	}

	for _, v := range b {
		if _, ok := hash[v]; ok {
			set = append(set, v)
		}
	}

	return set
}

func containsGeneric[T comparable](b []T, e T) bool {
	for _, v := range b {
		if v == e {
			return true
		}
	}
	return false
}
