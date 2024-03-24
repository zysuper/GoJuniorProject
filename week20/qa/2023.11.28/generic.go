package _023_11_28

func FindAll[T any](vals []T, filter func(t T) bool) []T {
	res := make([]T, 0, len(vals))
	for _, val := range vals {
		if filter(val) {
			res = append(res, val)
		}
	}
	return res
}
