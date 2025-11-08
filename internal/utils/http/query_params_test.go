package http

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestGetQueryParamInt(t *testing.T) {
	t.Run("returns parsed int", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users/42", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "42")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		value, err := GetQueryParamInt(req, "id")

		require.NoError(t, err)
		require.Equal(t, int64(42), value)
	})

	t.Run("returns error when value invalid", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users/foo", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "foo")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		_, err := GetQueryParamInt(req, "id")

		require.Error(t, err)
	})
}
