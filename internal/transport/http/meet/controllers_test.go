package meet

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	domainevent "billsplitter-monolith/internal/domain/event"
	vo "billsplitter-monolith/internal/domain/valueobject"
	appErrors "billsplitter-monolith/internal/errors"
	eventucmock "billsplitter-monolith/internal/mocks/usecase/event"
	usecaseevent "billsplitter-monolith/internal/usecase/event"
	"billsplitter-monolith/internal/utils"
	hu "billsplitter-monolith/internal/utils/http"
	"billsplitter-monolith/internal/utils/testkit"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errMeetTest = errors.New("meet-test-error")

func TestMeetController_CreateMeet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		eventUC := eventucmock.NewMockUseCase(t)
		ctrl := &controllerImpl{
			eventUC: eventUC,
			logger:  testkit.NewTestLogger(),
		}
		body := bytes.NewBufferString(`{"name":"Party","members":["Bob"]}`)
		req := httptest.NewRequest(http.MethodPost, "/meets", body)
		req = testkit.WithUserSession(req, vo.UserID(10))
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
		ctrl := &controllerImpl{eventUC: eventUC, logger: testkit.NewTestLogger()}
		req := httptest.NewRequest(http.MethodPost, "/meets", bytes.NewBufferString(`invalid`))
		req = testkit.WithUserSession(req, vo.UserID(1))
		rec := httptest.NewRecorder()

		ctrl.CreateMeet(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		eventUC := eventucmock.NewMockUseCase(t)
		ctrl := &controllerImpl{eventUC: eventUC, logger: testkit.NewTestLogger()}
		req := httptest.NewRequest(http.MethodPost, "/meets", bytes.NewBufferString(`{"name":"Party"}`))
		req = testkit.WithUserSession(req, vo.UserID(1))
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
		ctrl := &controllerImpl{eventUC: eventUC, logger: testkit.NewTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/meets", nil)
		req = testkit.WithUserSession(req, vo.UserID(7))
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
		ctrl := &controllerImpl{eventUC: eventucmock.NewMockUseCase(t), logger: testkit.NewTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/meets", nil)
		rec := httptest.NewRecorder()

		ctrl.FetchUserMeets(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("nil session", func(t *testing.T) {
		ctrl := &controllerImpl{eventUC: eventucmock.NewMockUseCase(t), logger: testkit.NewTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/meets", nil)
		req = testkit.WithSession(req, nil)
		rec := httptest.NewRecorder()

		ctrl.FetchUserMeets(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		eventUC := eventucmock.NewMockUseCase(t)
		ctrl := &controllerImpl{eventUC: eventUC, logger: testkit.NewTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/meets", nil)
		req = testkit.WithUserSession(req, vo.UserID(7))
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
		ctrl := &controllerImpl{eventUC: eventUC, logger: testkit.NewTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/meets/1", nil)
		req = testkit.WithRouteParams(req, map[string]string{"id": "1"})
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
		ctrl := &controllerImpl{eventUC: eventucmock.NewMockUseCase(t), logger: testkit.NewTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/meets/foo", nil)
		req = testkit.WithRouteParams(req, map[string]string{"id": "foo"})
		rec := httptest.NewRecorder()

		ctrl.GetMeetDetailsByID(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		eventUC := eventucmock.NewMockUseCase(t)
		ctrl := &controllerImpl{eventUC: eventUC, logger: testkit.NewTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/meets/1", nil)
		req = testkit.WithRouteParams(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		eventUC.EXPECT().
			GetMeetByID(mock.Anything, int64(1)).
			Return(nil, errMeetTest)

		ctrl.GetMeetDetailsByID(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("not found", func(t *testing.T) {
		eventUC := eventucmock.NewMockUseCase(t)
		ctrl := &controllerImpl{eventUC: eventUC, logger: testkit.NewTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/meets/1", nil)
		req = testkit.WithRouteParams(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		eventUC.EXPECT().
			GetMeetByID(mock.Anything, int64(1)).
			Return(nil, nil)

		ctrl.GetMeetDetailsByID(rec, req)

		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestMeetController_GetSummary(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		eventUC := eventucmock.NewMockUseCase(t)
		ctrl := &controllerImpl{eventUC: eventUC, logger: testkit.NewTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/meets/1/summary", nil)
		req = testkit.WithRouteParams(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		summary := &usecaseevent.EventSummary{
			Balances: []usecaseevent.Balance{
				{
					MemberID:   1,
					UserID:     utils.Ptr[vo.UserID](vo.UserID(1)),
					Name:       "Андрей",
					TotalPaid:  60,
					TotalShare: 50,
					Balance:    10,
				},
				{
					MemberID:   2,
					UserID:     utils.Ptr[vo.UserID](vo.UserID(2)),
					Name:       "Настя",
					TotalPaid:  0,
					TotalShare: 50,
					Balance:    -50,
				},
			},
			Settlements: []usecaseevent.Settlement{
				{FromMemberID: 2, ToMemberID: 1, Amount: 40},
			},
		}

		eventUC.EXPECT().
			CalculateSummary(mock.Anything, int64(1)).
			Return(summary, nil)

		ctrl.GetSummary(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		var resp EventSummary
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, summary.Balances[0].MemberID, resp.Balances[0].MemberID)
		require.Equal(t, utils.UserIDToInt64(summary.Balances[0].UserID), resp.Balances[0].UserID)
		require.Equal(t, summary.Settlements[0].Amount, resp.Settlements[0].Amount)
	})

	t.Run("invalid id", func(t *testing.T) {
		ctrl := &controllerImpl{eventUC: eventucmock.NewMockUseCase(t), logger: testkit.NewTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/meets/foo/summary", nil)
		req = testkit.WithRouteParams(req, map[string]string{"id": "foo"})
		rec := httptest.NewRecorder()

		ctrl.GetSummary(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		eventUC := eventucmock.NewMockUseCase(t)
		ctrl := &controllerImpl{eventUC: eventUC, logger: testkit.NewTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/meets/1/summary", nil)
		req = testkit.WithRouteParams(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		eventUC.EXPECT().
			CalculateSummary(mock.Anything, int64(1)).
			Return(nil, errMeetTest)

		ctrl.GetSummary(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("not found", func(t *testing.T) {
		eventUC := eventucmock.NewMockUseCase(t)
		ctrl := &controllerImpl{eventUC: eventUC, logger: testkit.NewTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/meets/1/summary", nil)
		req = testkit.WithRouteParams(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		eventUC.EXPECT().
			CalculateSummary(mock.Anything, int64(1)).
			Return(nil, appErrors.ErrEventNotFound)

		ctrl.GetSummary(rec, req)

		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}
