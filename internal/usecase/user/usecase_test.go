package user

import (
	"context"
	stderrors "errors"
	"testing"

	domainsession "billsplitter-monolith/internal/domain/session"
	domainuser "billsplitter-monolith/internal/domain/user"
	vo "billsplitter-monolith/internal/domain/valueobject"
	appErrors "billsplitter-monolith/internal/errors"
	sessionmock "billsplitter-monolith/internal/mocks/domain/session"
	usermock "billsplitter-monolith/internal/mocks/domain/user"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUseCase_GetByTgIDOrCreate(t *testing.T) {
	t.Run("user already exists", func(t *testing.T) {
		userSvc := usermock.NewMockService(t)
		sessionSvc := sessionmock.NewMockService(t)
		uc := New(userSvc, sessionSvc)
		ctx := context.Background()
		input := domainuser.User{Extra: domainuser.ExtraInfo{TelegramID: 123}}
		existing := &domainuser.User{ID: vo.UserID(1)}

		userSvc.EXPECT().
			GetByTelegramID(mock.Anything, input.Extra.TelegramID).
			Return(existing, nil)

		result, err := uc.GetByTgIDOrCreate(ctx, input)

		require.NoError(t, err)
		require.Equal(t, existing, result)
	})

	t.Run("create new user", func(t *testing.T) {
		userSvc := usermock.NewMockService(t)
		sessionSvc := sessionmock.NewMockService(t)
		uc := New(userSvc, sessionSvc)
		ctx := context.Background()
		input := domainuser.User{Username: "new", Extra: domainuser.ExtraInfo{TelegramID: 123}}
		created := &domainuser.User{ID: vo.UserID(2)}

		userSvc.EXPECT().
			GetByTelegramID(mock.Anything, input.Extra.TelegramID).
			Return(nil, appErrors.ErrUserNotFound)

		userSvc.EXPECT().
			Create(mock.Anything, input).
			Return(created, nil)

		result, err := uc.GetByTgIDOrCreate(ctx, input)

		require.NoError(t, err)
		require.Equal(t, created, result)
	})

	t.Run("unexpected error", func(t *testing.T) {
		userSvc := usermock.NewMockService(t)
		sessionSvc := sessionmock.NewMockService(t)
		uc := New(userSvc, sessionSvc)
		ctx := context.Background()
		input := domainuser.User{Extra: domainuser.ExtraInfo{TelegramID: 123}}
		svcErr := stderrors.New("boom")

		userSvc.EXPECT().
			GetByTelegramID(mock.Anything, input.Extra.TelegramID).
			Return(nil, svcErr)

		result, err := uc.GetByTgIDOrCreate(ctx, input)

		require.Nil(t, result)
		require.ErrorIs(t, err, svcErr)
	})

	t.Run("create error", func(t *testing.T) {
		userSvc := usermock.NewMockService(t)
		sessionSvc := sessionmock.NewMockService(t)
		uc := New(userSvc, sessionSvc)
		ctx := context.Background()
		input := domainuser.User{Extra: domainuser.ExtraInfo{TelegramID: 123}}
		svcErr := stderrors.New("create failed")

		userSvc.EXPECT().
			GetByTelegramID(mock.Anything, input.Extra.TelegramID).
			Return(nil, appErrors.ErrUserNotFound)

		userSvc.EXPECT().
			Create(mock.Anything, input).
			Return((*domainuser.User)(nil), svcErr)

		result, err := uc.GetByTgIDOrCreate(ctx, input)

		require.Nil(t, result)
		require.ErrorIs(t, err, svcErr)
	})
}

func TestUseCase_GetBySessionID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userSvc := usermock.NewMockService(t)
		sessionSvc := sessionmock.NewMockService(t)
		uc := New(userSvc, sessionSvc)
		ctx := context.Background()
		sessionID := "session"
		sessionUserID := vo.UserID(5)
		sess := &domainsession.Session{ID: sessionID, UserID: sessionUserID}
		userInfo := &domainuser.User{ID: sessionUserID}

		sessionSvc.EXPECT().
			GetByID(mock.Anything, sessionID).
			Return(sess, nil)

		userSvc.EXPECT().
			GetByID(mock.Anything, sessionUserID).
			Return(userInfo, nil)

		result, err := uc.GetBySessionID(ctx, sessionID)

		require.NoError(t, err)
		require.Equal(t, userInfo, result)
	})

	t.Run("session service error", func(t *testing.T) {
		userSvc := usermock.NewMockService(t)
		sessionSvc := sessionmock.NewMockService(t)
		uc := New(userSvc, sessionSvc)
		ctx := context.Background()
		sessionID := "session"
		svcErr := stderrors.New("session failure")

		sessionSvc.EXPECT().
			GetByID(mock.Anything, sessionID).
			Return((*domainsession.Session)(nil), svcErr)

		result, err := uc.GetBySessionID(ctx, sessionID)

		require.Nil(t, result)
		require.ErrorIs(t, err, svcErr)
	})

	t.Run("user service error", func(t *testing.T) {
		userSvc := usermock.NewMockService(t)
		sessionSvc := sessionmock.NewMockService(t)
		uc := New(userSvc, sessionSvc)
		ctx := context.Background()
		sessionID := "session"
		sessionUserID := vo.UserID(5)
		sess := &domainsession.Session{ID: sessionID, UserID: sessionUserID}
		svcErr := stderrors.New("user failure")

		sessionSvc.EXPECT().
			GetByID(mock.Anything, sessionID).
			Return(sess, nil)

		userSvc.EXPECT().
			GetByID(mock.Anything, sessionUserID).
			Return((*domainuser.User)(nil), svcErr)

		result, err := uc.GetBySessionID(ctx, sessionID)

		require.Nil(t, result)
		require.ErrorIs(t, err, svcErr)
	})
}
