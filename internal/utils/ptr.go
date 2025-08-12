package utils

func Ptr[T any](v T) *T {
	return &v
}

func SafeDereference[T any](v *T) T {
	var defaultValue T
	if v == nil {
		return defaultValue
	}
	return *v
}
