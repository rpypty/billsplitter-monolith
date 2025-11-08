package meet

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	domainevent "billsplitter-monolith/internal/domain/event"
	domainsession "billsplitter-monolith/internal/domain/session"
	vo "billsplitter-monolith/internal/domain/valueobject"
	eventucmock "billsplitter-monolith/internal/mocks/usecase/event"
	"billsplitter-monolith/internal/transport/http/middleware"
	usecaseevent "billsplitter-monolith/internal/usecase/event"
	hu "billsplitter-monolith/internal/utils/http"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMeetController_CreateMeet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		eventUC := eventucmock.NewMockUseCase(t)
		ctrl := &controllerImpl{
			eventUC: eventUC,
			logger:  newTestLogger(),
		}
		body := bytes.NewBufferString(`{"name":"Party","members":["Bob"]}`)
		req := httptest.NewRequest(http.MethodPost, "/meets", body)
		req = withMeetSession(req, vo.UserID(10))
		rec := httptest.NewRecorder()

		eventUC.EXPECT().
			CreateMeet(mock.Anything, mock.MatchedBy(func(r usecaseevent.CreateMeetRq) bool {
				return r.EventName == "Party" && r.CreatedByUserID == vo.UserID(10)
			})).
			Return(int64(55), nil)

		ctrl.CreateMeet(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		var resp hu.ResponseID
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, int64(55), resp.ID)
	})

	t.Run("decode error", func(t *testing.T) {
		eventUC := eventucmock.NewMockUseCase(t)
		ctrl := &controllerImpl{eventUC: eventUC, logger: newTestLogger()}
		req := httptest.NewRequest(http.MethodPost, "/meets", bytes.NewBufferString(`invalid`))
		req = withMeetSession(req, vo.UserID(1))
		rec := httptest.NewRecorder()

		ctrl.CreateMeet(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		eventUC := eventucmock.NewMockUseCase(t)
		ctrl := &controllerImpl{eventUC: eventUC, logger: newTestLogger()}
		req := httptest.NewRequest(http.MethodPost, "/meets", bytes.NewBufferString(`{"name":"Party"}`))
		req = withMeetSession(req, vo.UserID(1))
		rec := httptest.NewRecorder()

		eventUC.EXPECT().
			CreateMeet(mock.Anything, mock.Anything).
			Return(int64(0), errMeetTest)

		ctrl.CreateMeet(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestMeetController_FetchUserMeets(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		eventUC := eventucmock.NewMockUseCase(t)
		ctrl := &controllerImpl{eventUC: eventUC, logger: newTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/meets", nil)
		req = withMeetSession(req, vo.UserID(7))
		rec := httptest.NewRecorder()

		events := []domainevent.Event{
			{ID: 1, Name: "Party", CreatedByUserID: vo.UserID(7)},
		}

		eventUC.EXPECT().
			FetchUserMeets(mock.Anything, vo.UserID(7)).
			Return(events, nil)

		ctrl.FetchUserMeets(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		var resp []Event
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Len(t, resp, 1)
		require.Equal(t, int64(1), resp[0].ID)
	})

	t.Run("missing session", func(t *testing.T) {
		ctrl := &controllerImpl{eventUC: eventucmock.NewMockUseCase(t), logger: newTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/meets", nil)
		rec := httptest.NewRecorder()

		ctrl.FetchUserMeets(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("nil session", func(t *testing.T) {
		ctrl := &controllerImpl{eventUC: eventucmock.NewMockUseCase(t), logger: newTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/meets", nil)
		req = req.WithContext(middleware.ContextWithSessionForTest(req.Context(), nil))
		rec := httptest.NewRecorder()

		ctrl.FetchUserMeets(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		eventUC := eventucmock.NewMockUseCase(t)
		ctrl := &controllerImpl{eventUC: eventUC, logger: newTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/meets", nil)
		req = withMeetSession(req, vo.UserID(7))
		rec := httptest.NewRecorder()

		eventUC.EXPECT().
			FetchUserMeets(mock.Anything, vo.UserID(7)).
			Return(nil, errMeetTest)

		ctrl.FetchUserMeets(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestMeetController_GetMeetDetailsByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		eventUC := eventucmock.NewMockUseCase(t)
		ctrl := &controllerImpl{eventUC: eventUC, logger: newTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/meets/1", nil)
		req = withRouteParams(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		eventObj := &domainevent.Event{ID: 1, Name: "Party"}

		eventUC.EXPECT().
			GetMeetByID(mock.Anything, int64(1)).
			Return(eventObj, nil)

		ctrl.GetMeetDetailsByID(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		var resp Event
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, int64(1), resp.ID)
	})

	t.Run("invalid id", func(t *testing.T) {
		ctrl := &controllerImpl{eventUC: eventucmock.NewMockUseCase(t), logger: newTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/meets/foo", nil)
		req = withRouteParams(req, map[string]string{"id": "foo"})
		rec := httptest.NewRecorder()

		ctrl.GetMeetDetailsByID(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		eventUC := eventucmock.NewMockUseCase(t)
		ctrl := &controllerImpl{eventUC: eventUC, logger: newTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/meets/1", nil)
		req = withRouteParams(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		eventUC.EXPECT().
			GetMeetByID(mock.Anything, int64(1)).
			Return(nil, errMeetTest)

		ctrl.GetMeetDetailsByID(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("not found", func(t *testing.T) {
		eventUC := eventucmock.NewMockUseCase(t)
		ctrl := &controllerImpl{eventUC: eventUC, logger: newTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/meets/1", nil)
		req = withRouteParams(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		eventUC.EXPECT().
			GetMeetByID(mock.Anything, int64(1)).
			Return(nil, nil)

		ctrl.GetMeetDetailsByID(rec, req)

		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func withMeetSession(req *http.Request, userID vo.UserID) *http.Request {
	sess := &domainsession.Session{UserID: userID}
	return req.WithContext(middleware.ContextWithSessionForTest(req.Context(), sess))
}

func withRouteParams(req *http.Request, params map[string]string) *http.Request {
	routeCtx := chi.NewRouteContext()
	for k, v := range params {
		routeCtx.URLParams.Add(k, v)
	}
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx)
	return req.WithContext(ctx)
}

var errMeetTest = errors.New("meet-test-error")

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
