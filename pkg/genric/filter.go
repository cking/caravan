package genric

func FilterArray[T any](slice []T, f func(T) bool) []T {
	var n []T

	for _, e := range slice {
		if f(e) {
			n = append(n, e)
		}
	}
	return n
}

func FilterMap[TKey comparable, TValue any](slice map[TKey]TValue, f func(TKey, TValue) bool) map[TKey]TValue {
	m := make(map[TKey]TValue)

	for k, v := range slice {
		if f(k, v) {
			m[k] = v
		}
	}
	return m
}
