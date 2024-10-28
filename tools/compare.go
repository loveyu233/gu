package tools

func Min[T int | int32 | int64 | float32 | float64](x, y T) T {
	if x > y {
		return y
	}
	return x
}

func Max[T int | int32 | int64 | float32 | float64](x, y T) T {
	if x > y {
		return x
	}
	return y
}
