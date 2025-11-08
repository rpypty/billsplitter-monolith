package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	domainsession "billsplitter-monolith/internal/domain/session"
	domainuser "billsplitter-monolith/internal/domain/user"
	vo "billsplitter-monolith/internal/domain/valueobject"
	sessionmock "billsplitter-monolith/internal/mocks/domain/session"
	usermock "billsplitter-monolith/internal/mocks/domain/user"
	userucmock "billsplitter-monolith/internal/mocks/usecase/user"
	"billsplitter-monolith/internal/transport/http/middleware"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestController_LoginTelegram(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userUC := userucmock.NewMockUseCase(t)
		userSvc := usermock.NewMockService(t)
		sessionSvc := sessionmock.NewMockService(t)
		ctrl := NewController(userUC, userSvc, sessionSvc, newTestLogger())

		body := bytes.NewBufferString(`{"username":"nick","firstName":"Nick","lastName":"Doe","telegramID":123}`)
		req := httptest.NewRequest(http.MethodPost, "/auth/login/telegram", body)
		rec := httptest.NewRecorder()

		expectedUser := &domainuser.User{ID: vo.UserID(1), Username: "nick"}

		userUC.EXPECT().
			GetByTgIDOrCreate(mock.Anything, mock.MatchedBy(func(u domainuser.User) bool {
				return u.Username == "nick" && u.Extra.TelegramID == 123
			})).
			Return(expectedUser, nil)

		sessionSvc.EXPECT().
			Create(mock.Anything, expectedUser.ID).
			Return("session-id", nil)

		ctrl.LoginTelegram(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		var resp LoginTelegramRes
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, "session-id", resp.SessionID)
		require.Equal(t, expectedUser, resp.UserInfo)
	})

	t.Run("invalid payload", func(t *testing.T) {
		userUC := userucmock.NewMockUseCase(t)
		sessionSvc := sessionmock.NewMockService(t)
		ctrl := NewController(userUC, nil, sessionSvc, newTestLogger())

		body := bytes.NewBufferString(`{"username":"nick"}`)
		req := httptest.NewRequest(http.MethodPost, "/auth/login/telegram", body)
		rec := httptest.NewRecorder()

		ctrl.LoginTelegram(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("user usecase error", func(t *testing.T) {
		userUC := userucmock.NewMockUseCase(t)
		sessionSvc := sessionmock.NewMockService(t)
		ctrl := NewController(userUC, nil, sessionSvc, newTestLogger())
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"username":"nick","telegramID":1}`))
		rec := httptest.NewRecorder()

		userUC.EXPECT().
			GetByTgIDOrCreate(mock.Anything, mock.Anything).
			Return(nil, errTest)

		ctrl.LoginTelegram(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("session create error", func(t *testing.T) {
		userUC := userucmock.NewMockUseCase(t)
		sessionSvc := sessionmock.NewMockService(t)
		ctrl := NewController(userUC, nil, sessionSvc, newTestLogger())
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"username":"nick","telegramID":1}`))
		rec := httptest.NewRecorder()

		expectedUser := &domainuser.User{ID: vo.UserID(1)}
		userUC.EXPECT().
			GetByTgIDOrCreate(mock.Anything, mock.Anything).
			Return(expectedUser, nil)

		sessionSvc.EXPECT().
			Create(mock.Anything, expectedUser.ID).
			Return("", errTest)

		ctrl.LoginTelegram(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestController_Me(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userSvc := usermock.NewMockService(t)
		sessionSvc := sessionmock.NewMockService(t)
		ctrl := NewController(nil, userSvc, sessionSvc, newTestLogger())
		req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
		sess := &domainsession.Session{UserID: vo.UserID(10)}
		req = req.WithContext(middleware.ContextWithSessionForTest(req.Context(), sess))
		rec := httptest.NewRecorder()

		expectedUser := &domainuser.User{ID: vo.UserID(10), Username: "nick"}

		userSvc.EXPECT().
			GetByID(mock.Anything, sess.UserID).
			Return(expectedUser, nil)

		ctrl.Me(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		var resp MeRes
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, expectedUser, resp.User)
	})

	t.Run("missing session", func(t *testing.T) {
		ctrl := NewController(nil, nil, nil, newTestLogger())
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)

		ctrl.Me(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("user service error", func(t *testing.T) {
		userSvc := usermock.NewMockService(t)
		ctrl := NewController(nil, userSvc, nil, newTestLogger())
		req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
		sess := &domainsession.Session{UserID: vo.UserID(5)}
		req = req.WithContext(middleware.ContextWithSessionForTest(req.Context(), sess))
		rec := httptest.NewRecorder()

		userSvc.EXPECT().
			GetByID(mock.Anything, sess.UserID).
			Return(nil, errTest)

		ctrl.Me(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

var errTest = errors.New("test-error")

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
