package maps

// A key-value pair
type KV[K, V any] struct {
	Key   K
	Value V
}

func Keys[K comparable, V any](m map[K]V) []K {
	items := make([]K, 0, len(m))
	for k, _ := range m {
		items = append(items, k)
	}
	return items
}

func Items[K comparable, V any](m map[K]V) []KV[K, V] {
	items := make([]KV[K, V], 0, len(m))
	for k, v := range m {
		items = append(items, KV[K, V]{
			Key:   k,
			Value: v,
		})
	}
	return items
}
