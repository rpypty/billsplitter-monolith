package impl

import (
	"context"
	stderrors "errors"
	"testing"

	domain "billsplitter-monolith/internal/domain/payment_method"
	vo "billsplitter-monolith/internal/domain/valueobject"
	paymentmethodmock "billsplitter-monolith/internal/mocks/domain/payment_method"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_FetchByUserID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := paymentmethodmock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		userID := vo.UserID(42)
		expected := []domain.PaymentMethod{
			{ID: 1, UserID: int64(userID), Name: "Mock Pay"},
		}

		repo.EXPECT().
			FetchByUserID(mock.Anything, userID).
			Return(expected, nil)

		result, err := svc.FetchByUserID(ctx, userID)

		require.NoError(t, err)
		require.Equal(t, expected, result)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := paymentmethodmock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		userID := vo.UserID(42)
		repoErr := stderrors.New("fetch failed")

		repo.EXPECT().
			FetchByUserID(mock.Anything, userID).
			Return(nil, repoErr)

		result, err := svc.FetchByUserID(ctx, userID)

		require.Nil(t, result)
		require.ErrorIs(t, err, repoErr)
	})
}

func TestService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := paymentmethodmock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		input := domain.PaymentMethod{Name: "Alpha", UserID: 1}

		repo.EXPECT().
			Create(mock.Anything, input).
			Return(input, nil)

		result, err := svc.Create(ctx, input)

		require.NoError(t, err)
		require.Equal(t, input, result)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := paymentmethodmock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		input := domain.PaymentMethod{Name: "Alpha", UserID: 1}
		repoErr := stderrors.New("create failed")

		repo.EXPECT().
			Create(mock.Anything, input).
			Return(domain.PaymentMethod{}, repoErr)

		_, err := svc.Create(ctx, input)

		require.ErrorIs(t, err, repoErr)
	})
}

func TestService_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := paymentmethodmock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		update := domain.PaymentMethod{Name: "New"}

		repo.EXPECT().
			Update(mock.Anything, int64(1), update).
			Return(nil)

		err := svc.Update(ctx, 1, update)

		require.NoError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := paymentmethodmock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		update := domain.PaymentMethod{Name: "New"}
		repoErr := stderrors.New("update failed")

		repo.EXPECT().
			Update(mock.Anything, int64(1), update).
			Return(repoErr)

		err := svc.Update(ctx, 1, update)

		require.ErrorIs(t, err, repoErr)
	})
}

func TestService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := paymentmethodmock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()

		repo.EXPECT().
			Delete(mock.Anything, int64(1)).
			Return(nil)

		err := svc.Delete(ctx, 1)

		require.NoError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := paymentmethodmock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		repoErr := stderrors.New("delete failed")

		repo.EXPECT().
			Delete(mock.Anything, int64(1)).
			Return(repoErr)

		err := svc.Delete(ctx, 1)

		require.ErrorIs(t, err, repoErr)
	})
}
