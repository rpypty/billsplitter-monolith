package http

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type sampleRequest struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age"`
}

func TestDecodeReq(t *testing.T) {
	t.Run("successful decode and validation", func(t *testing.T) {
		body := bytes.NewBufferString(`{"name":"Alice","age":30}`)
		req := httptest.NewRequest("POST", "/test", body)

		result, err := DecodeReq[sampleRequest](req)

		require.NoError(t, err)
		require.Equal(t, &sampleRequest{Name: "Alice", Age: 30}, result)
	})

	t.Run("validation failure", func(t *testing.T) {
		body := bytes.NewBufferString(`{"age":30}`)
		req := httptest.NewRequest("POST", "/test", body)

		result, err := DecodeReq[sampleRequest](req)

		require.Nil(t, result)
		require.Error(t, err)
	})
}

func TestMustMarshal(t *testing.T) {
	type sample struct {
		Value string `json:"value"`
	}

	bytes := MustMarshal(sample{Value: "ok"})

	require.JSONEq(t, `{"value":"ok"}`, string(bytes))
}
