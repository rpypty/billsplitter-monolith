package utils

import (
	"testing"

	vo "billsplitter-monolith/internal/domain/valueobject"

	"github.com/stretchr/testify/require"
)

func TestPtr(t *testing.T) {
	t.Run("returns pointer to value", func(t *testing.T) {
		val := "hello"
		ptr := Ptr(val)

		require.NotNil(t, ptr)
		require.Equal(t, val, *ptr)
		require.NotSame(t, &val, ptr)
	})
}

func TestSafeDereference(t *testing.T) {
	t.Run("returns zero value for nil pointer", func(t *testing.T) {
		var ptr *int

		val := SafeDereference(ptr)

		require.Zero(t, val)
	})

	t.Run("returns pointed value", func(t *testing.T) {
		val := 42

		result := SafeDereference(&val)

		require.Equal(t, val, result)
	})
}

func TestUserIDToInt64(t *testing.T) {
	t.Run("nil input returns nil", func(t *testing.T) {
		require.Nil(t, UserIDToInt64(nil))
	})

	t.Run("converts value", func(t *testing.T) {
		id := vo.UserID(55)

		result := UserIDToInt64(&id)

		require.NotNil(t, result)
		require.Equal(t, int64(id), *result)
	})
}
