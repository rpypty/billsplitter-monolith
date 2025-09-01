package utils

import vo "billsplitter-monolith/internal/domain/valueobject"

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

func UserIDToInt64(v *vo.UserID) *int64 {
	if v == nil {
		return nil
	}

	return Ptr(int64(*v))
}
