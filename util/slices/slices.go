package slices

import (
	"sort"

	"golang.org/x/exp/constraints"
)

func SortOn[K constraints.Ordered, T any](items []T, key func(T) K) {
	sort.Slice(items, func(i, j int) bool {
		return key(items[i]) < key(items[j])
	})
}
