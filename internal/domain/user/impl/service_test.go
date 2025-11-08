package impl

import (
	"context"
	stderrors "errors"
	"testing"

	domain "billsplitter-monolith/internal/domain/user"
	vo "billsplitter-monolith/internal/domain/valueobject"
	appErrors "billsplitter-monolith/internal/errors"
	usermock "billsplitter-monolith/internal/mocks/domain/user"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_Fetch(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := usermock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		filter := domain.FetchFilter{UserIDs: []int64{1, 2}}
		expected := []domain.User{{ID: vo.UserID(1), Username: "tester"}}

		repo.EXPECT().
			Fetch(mock.Anything, filter).
			Return(expected, nil)

		users, err := svc.Fetch(ctx, filter)

		require.NoError(t, err)
		require.Equal(t, expected, users)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := usermock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		repoErr := stderrors.New("fetch failed")

		repo.EXPECT().
			Fetch(mock.Anything, mock.Anything).
			Return(nil, repoErr)

		users, err := svc.Fetch(ctx, domain.FetchFilter{})

		require.Nil(t, users)
		require.ErrorContains(t, err, "ServiceImpl.Fetch")
		require.ErrorIs(t, err, repoErr)
	})
}

func TestService_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := usermock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		userID := vo.UserID(10)
		expected := &domain.User{ID: userID, Username: "tester"}

		repo.EXPECT().
			GetByID(mock.Anything, userID).
			Return(expected, nil)

		user, err := svc.GetByID(ctx, userID)

		require.NoError(t, err)
		require.Equal(t, expected, user)
	})

	t.Run("user not found", func(t *testing.T) {
		repo := usermock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		userID := vo.UserID(7)

		repo.EXPECT().
			GetByID(mock.Anything, userID).
			Return(nil, nil)

		user, err := svc.GetByID(ctx, userID)

		require.Nil(t, user)
		require.ErrorContains(t, err, "user not found")
	})

	t.Run("repository error", func(t *testing.T) {
		repo := usermock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		userID := vo.UserID(7)
		repoErr := stderrors.New("boom")

		repo.EXPECT().
			GetByID(mock.Anything, userID).
			Return(nil, repoErr)

		_, err := svc.GetByID(ctx, userID)

		require.ErrorContains(t, err, "ServiceImpl.GetById")
		require.ErrorIs(t, err, repoErr)
	})
}

func TestService_GetByTelegramID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := usermock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		telegramID := int64(1000)
		expected := &domain.User{
			ID:       vo.UserID(1),
			Username: "telegram-user",
			Extra: domain.ExtraInfo{
				TelegramID: telegramID,
			},
		}

		repo.EXPECT().
			GetByTelegramID(mock.Anything, telegramID).
			Return(expected, nil)

		user, err := svc.GetByTelegramID(ctx, telegramID)

		require.NoError(t, err)
		require.Equal(t, expected, user)
	})

	// repo returns nil, nil
	t.Run("user not found", func(t *testing.T) {
		repo := usermock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		telegramID := int64(1000)

		repo.EXPECT().
			GetByTelegramID(mock.Anything, telegramID).
			Return(nil, nil)

		_, err := svc.GetByTelegramID(ctx, telegramID)

		require.ErrorContains(t, err, "user not found")
		require.ErrorIs(t, err, appErrors.ErrUserNotFound)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := usermock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		telegramID := int64(1000)
		repoErr := stderrors.New("db err")

		repo.EXPECT().
			GetByTelegramID(mock.Anything, telegramID).
			Return(nil, repoErr)

		_, err := svc.GetByTelegramID(ctx, telegramID)

		require.ErrorContains(t, err, "ServiceImpl.GetByTelegramID")
		require.ErrorIs(t, err, repoErr)
	})
}

func TestService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := usermock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		userInput := validUser()
		created := userInput
		created.ID = vo.UserID(99)

		repo.EXPECT().
			GetByTelegramID(mock.Anything, userInput.Extra.TelegramID).
			Return(nil, nil)
		repo.EXPECT().
			Create(mock.Anything, userInput).
			Return(&created, nil)

		user, err := svc.Create(ctx, userInput)

		require.NoError(t, err)
		require.Equal(t, &created, user)
	})

	t.Run("validation error", func(t *testing.T) {
		repo := usermock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		invalid := domain.User{Username: "abc"}

		user, err := svc.Create(ctx, invalid)

		require.Nil(t, user)
		require.ErrorContains(t, err, "validation error")
	})

	t.Run("telegram already exists", func(t *testing.T) {
		repo := usermock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		userInput := validUser()

		repo.EXPECT().
			GetByTelegramID(mock.Anything, userInput.Extra.TelegramID).
			Return(&domain.User{}, nil)

		user, err := svc.Create(ctx, userInput)

		require.Nil(t, user)
		require.ErrorContains(t, err, "already exists")
	})

	t.Run("getByTelegram error", func(t *testing.T) {
		repo := usermock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		userInput := validUser()
		repoErr := stderrors.New("lookup failed")

		repo.EXPECT().
			GetByTelegramID(mock.Anything, userInput.Extra.TelegramID).
			Return(nil, repoErr)

		_, err := svc.Create(ctx, userInput)

		require.ErrorContains(t, err, "ServiceImpl.Create")
		require.ErrorIs(t, err, repoErr)
	})

	t.Run("create error", func(t *testing.T) {
		repo := usermock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		userInput := validUser()
		repoErr := stderrors.New("insert failed")

		repo.EXPECT().
			GetByTelegramID(mock.Anything, userInput.Extra.TelegramID).
			Return(nil, nil)
		repo.EXPECT().
			Create(mock.Anything, userInput).
			Return(nil, repoErr)

		_, err := svc.Create(ctx, userInput)

		require.ErrorContains(t, err, "ServiceImpl.Create")
		require.ErrorIs(t, err, repoErr)
	})
}

func TestService_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := usermock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		userInput := validUser()
		userID := vo.UserID(1)

		repo.EXPECT().
			Update(mock.Anything, userID, userInput).
			Return(nil)

		err := svc.Update(ctx, userID, userInput)

		require.NoError(t, err)
	})

	t.Run("validation error", func(t *testing.T) {
		repo := usermock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		userID := vo.UserID(1)
		invalid := domain.User{Username: "abc"}

		err := svc.Update(ctx, userID, invalid)

		require.ErrorContains(t, err, "validation error")
	})

	t.Run("repository error", func(t *testing.T) {
		repo := usermock.NewMockRepository(t)
		svc := New(repo)
		ctx := context.Background()
		userInput := validUser()
		userID := vo.UserID(1)
		repoErr := stderrors.New("update failed")

		repo.EXPECT().
			Update(mock.Anything, userID, userInput).
			Return(repoErr)

		err := svc.Update(ctx, userID, userInput)

		require.ErrorContains(t, err, "ServiceImpl.Update")
		require.ErrorIs(t, err, repoErr)
	})
}

func validUser() domain.User {
	return domain.User{
		Username:  "tester",
		FirstName: "Test",
		LastName:  "User",
		Extra: domain.ExtraInfo{
			TelegramID: 12345,
		},
	}
}
