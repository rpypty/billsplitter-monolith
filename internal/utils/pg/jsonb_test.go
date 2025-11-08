package pg

import (
	"testing"

	"github.com/jackc/pgtype"
	"github.com/stretchr/testify/require"
)

type samplePayload struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestToJsonb(t *testing.T) {
	t.Run("converts struct to jsonb", func(t *testing.T) {
		payload := &samplePayload{Name: "Alice", Age: 30}

		jsonb, err := ToJsonb(payload)

		require.NoError(t, err)
		require.NotNil(t, jsonb)

		var out samplePayload
		require.NoError(t, jsonb.AssignTo(&out))
		require.Equal(t, *payload, out)
	})

	t.Run("nil payload returns empty object", func(t *testing.T) {
		var payload *samplePayload

		jsonb, err := ToJsonb(payload)

		require.NoError(t, err)
		require.NotNil(t, jsonb)

		var out map[string]any
		require.NoError(t, jsonb.AssignTo(&out))
		require.Equal(t, map[string]any{}, out)
	})
}

func TestFromJsonb(t *testing.T) {
	t.Run("nil jsonb returns nil", func(t *testing.T) {
		value, err := FromJsonb[samplePayload](nil)

		require.NoError(t, err)
		require.Nil(t, value)
	})

	t.Run("converts jsonb to struct", func(t *testing.T) {
		payload := samplePayload{Name: "Bob", Age: 40}
		jsonb := pgtype.JSONB{}
		require.NoError(t, jsonb.Set(payload))

		value, err := FromJsonb[samplePayload](&jsonb)

		require.NoError(t, err)
		require.NotNil(t, value)
		require.Equal(t, payload, *value)
	})
}
