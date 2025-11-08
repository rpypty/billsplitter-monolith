package session

import (
	"context"
	"testing"
	"time"

	domain "billsplitter-monolith/internal/domain/session"

	"github.com/stretchr/testify/require"
)

func TestCache_Set(t *testing.T) {
	t.Run("stores session with ttl", func(t *testing.T) {
		cache := NewMemCache()
		ctx := context.Background()
		sess := &domain.Session{ID: "abc", UserID: 1}

		err := cache.Set(ctx, sess.ID, sess, time.Minute)

		require.NoError(t, err)

		stored, err := cache.Get(ctx, sess.ID)

		require.NoError(t, err)
		require.Equal(t, sess, stored)
	})

	t.Run("overwrites existing entry", func(t *testing.T) {
		cache := NewMemCache()
		ctx := context.Background()
		first := &domain.Session{ID: "same", UserID: 1}
		second := &domain.Session{ID: "same", UserID: 2}

		require.NoError(t, cache.Set(ctx, first.ID, first, time.Minute))
		require.NoError(t, cache.Set(ctx, second.ID, second, time.Minute))

		stored, err := cache.Get(ctx, first.ID)

		require.NoError(t, err)
		require.Equal(t, second, stored)
	})
}

func TestCache_Get(t *testing.T) {
	t.Run("returns session when present", func(t *testing.T) {
		cache := NewMemCache()
		ctx := context.Background()
		sess := &domain.Session{ID: "abc", UserID: 1}
		require.NoError(t, cache.Set(ctx, sess.ID, sess, time.Minute))

		stored, err := cache.Get(ctx, sess.ID)

		require.NoError(t, err)
		require.Equal(t, sess, stored)
	})

	t.Run("returns nil when not found", func(t *testing.T) {
		cache := NewMemCache()
		ctx := context.Background()

		stored, err := cache.Get(ctx, "missing")

		require.NoError(t, err)
		require.Nil(t, stored)
	})

	t.Run("returns nil when expired", func(t *testing.T) {
		cache := NewMemCache()
		ctx := context.Background()
		sess := &domain.Session{ID: "abc", UserID: 1}
		require.NoError(t, cache.Set(ctx, sess.ID, sess, -time.Second))

		stored, err := cache.Get(ctx, sess.ID)

		require.NoError(t, err)
		require.Nil(t, stored)
	})
}
