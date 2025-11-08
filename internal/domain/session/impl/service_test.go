package impl

import (
	"context"
	stderrors "errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"billsplitter-monolith/internal/domain/session"
	vo "billsplitter-monolith/internal/domain/valueobject"
	appErrors "billsplitter-monolith/internal/errors"
	sessionmock "billsplitter-monolith/internal/mocks/domain/session"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := sessionmock.NewMockRepository(t)
		svc := New(repo, newTestLogger())
		ctx := context.Background()
		sessionID := "session-id"
		expire := time.Now().Add(time.Minute)
		expected := &session.Session{
			ID:       sessionID,
			UserID:   vo.UserID(1),
			ExpireAt: &expire,
		}

		repo.EXPECT().
			GetByID(mock.Anything, sessionID).
			Return(expected, nil)

		obj, err := svc.GetByID(ctx, sessionID)

		require.NoError(t, err)
		require.Same(t, expected, obj)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := sessionmock.NewMockRepository(t)
		svc := New(repo, newTestLogger())
		ctx := context.Background()
		sessionID := "session-id"
		repoErr := stderrors.New("db error")

		repo.EXPECT().
			GetByID(mock.Anything, sessionID).
			Return(nil, repoErr)

		_, err := svc.GetByID(ctx, sessionID)

		require.ErrorContains(t, err, "ServiceImpl.GetByID")
		require.ErrorIs(t, err, repoErr)
	})

	t.Run("not found", func(t *testing.T) {
		repo := sessionmock.NewMockRepository(t)
		svc := New(repo, newTestLogger())
		ctx := context.Background()
		sessionID := "session-id"

		repo.EXPECT().
			GetByID(mock.Anything, sessionID).
			Return(nil, nil)

		_, err := svc.GetByID(ctx, sessionID)

		require.ErrorContains(t, err, "ServiceImpl.GetByID")
		require.ErrorIs(t, err, appErrors.ErrSessionNotFound)
	})

	t.Run("session expired", func(t *testing.T) {
		repo := sessionmock.NewMockRepository(t)
		svc := New(repo, newTestLogger())
		ctx := context.Background()
		sessionID := "session-id"
		expire := time.Now().Add(-time.Minute)

		repo.EXPECT().
			GetByID(mock.Anything, sessionID).
			Return(&session.Session{
				ID:       sessionID,
				UserID:   vo.UserID(1),
				ExpireAt: &expire,
			}, nil)

		_, err := svc.GetByID(ctx, sessionID)

		require.ErrorContains(t, err, "ServiceImpl.GetByID")
		require.ErrorIs(t, err, appErrors.ErrSessionExpired)
	})
}

func TestService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := sessionmock.NewMockRepository(t)
		svc := New(repo, newTestLogger())
		ctx := context.Background()
		userID := vo.UserID(42)
		start := time.Now()

		repo.EXPECT().
			Create(mock.Anything, mock.MatchedBy(func(s session.Session) bool {
				require.Equal(t, userID, s.UserID)
				require.NotNil(t, s.ExpireAt)
				require.WithinDuration(t, start.Add(sessionLiveTime), *s.ExpireAt, 2*time.Second)
				require.NotEmpty(t, s.ID)
				return true
			})).
			RunAndReturn(func(ctx context.Context, s session.Session) (*session.Session, error) {
				return &s, nil
			})

		sessionID, err := svc.Create(ctx, userID)

		require.NoError(t, err)
		require.NotEmpty(t, sessionID)
	})

	t.Run("validation error", func(t *testing.T) {
		repo := sessionmock.NewMockRepository(t)
		svc := New(repo, newTestLogger())
		ctx := context.Background()

		sessionID, err := svc.Create(ctx, 0)

		require.Empty(t, sessionID)
		require.ErrorContains(t, err, "validation error")
	})

	t.Run("repository error", func(t *testing.T) {
		repo := sessionmock.NewMockRepository(t)
		svc := New(repo, newTestLogger())
		ctx := context.Background()
		userID := vo.UserID(42)
		repoErr := stderrors.New("insert failed")

		repo.EXPECT().
			Create(mock.Anything, mock.Anything).
			Return(nil, repoErr)

		_, err := svc.Create(ctx, userID)

		require.ErrorContains(t, err, "ServiceImpl.Create")
		require.ErrorIs(t, err, repoErr)
	})
}

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
