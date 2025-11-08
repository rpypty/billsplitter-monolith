package utils

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewUUIDv7(t *testing.T) {
	t.Run("returns valid uuid string", func(t *testing.T) {
		first := NewUUIDv7()
		second := NewUUIDv7()

		require.NotEmpty(t, first)
		require.NotEmpty(t, second)

		parsedFirst, err := uuid.Parse(first)
		require.NoError(t, err)
		require.Equal(t, uuid.Version(7), parsedFirst.Version())

		parsedSecond, err := uuid.Parse(second)
		require.NoError(t, err)

		require.NotEqual(t, first, second, "should generate unique ids")
		require.NotEqual(t, parsedFirst, parsedSecond)
	})
}
