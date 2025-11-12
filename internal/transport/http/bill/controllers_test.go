package bill

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	domainbill "billsplitter-monolith/internal/domain/bill"
	vo "billsplitter-monolith/internal/domain/valueobject"
	apperrors "billsplitter-monolith/internal/errors"
	billusecasemock "billsplitter-monolith/internal/mocks/usecase/bill"
	billuc "billsplitter-monolith/internal/usecase/bill"
	hu "billsplitter-monolith/internal/utils/http"
	"billsplitter-monolith/internal/utils/testkit"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errControllerTest = errors.New("bill-controller-test-error")

func TestBillController_CreateBill(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		uc := billusecasemock.NewMockUseCase(t)
		ctrl := &controllerImpl{billUC: uc, logger: testkit.NewTestLogger()}

		body := bytes.NewBufferString(`{"event_id":1,"name":"Dinner","total_amount":1000,"currency":"BYN","split_type":"even","participants":[{"member_id":10,"amount":500}]}`)
		req := httptest.NewRequest(http.MethodPost, "/bills", body)
		req = testkit.WithUserSession(req, vo.UserID(5))
		rec := httptest.NewRecorder()

		uc.EXPECT().
			CreateBill(mock.Anything, mock.MatchedBy(func(r billuc.CreateBillRq) bool {
				return r.EventID == 1 && r.CreatedBy == vo.UserID(5) && len(r.Participants) == 1 && r.Participants[0].MemberID == 10
			})).
			Return(int64(77), nil)

		ctrl.CreateBill(rec, req)

		require.Equal(t, http.StatusCreated, rec.Code)
		var resp hu.ResponseID
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, int64(77), resp.ID)
	})

	t.Run("decode error", func(t *testing.T) {
		ctrl := &controllerImpl{billUC: billusecasemock.NewMockUseCase(t), logger: testkit.NewTestLogger()}
		req := httptest.NewRequest(http.MethodPost, "/bills", bytes.NewBufferString("invalid"))
		req = testkit.WithUserSession(req, vo.UserID(1))
		rec := httptest.NewRecorder()

		ctrl.CreateBill(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("missing session", func(t *testing.T) {
		ctrl := &controllerImpl{billUC: billusecasemock.NewMockUseCase(t), logger: testkit.NewTestLogger()}
		req := httptest.NewRequest(http.MethodPost, "/bills", nil)
		rec := httptest.NewRecorder()

		ctrl.CreateBill(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("usecase event not found", func(t *testing.T) {
		uc := billusecasemock.NewMockUseCase(t)
		ctrl := &controllerImpl{billUC: uc, logger: testkit.NewTestLogger()}
		req := httptest.NewRequest(http.MethodPost, "/bills", bytes.NewBufferString(`{"event_id":1}`))
		req = testkit.WithUserSession(req, vo.UserID(1))
		rec := httptest.NewRecorder()

		uc.EXPECT().CreateBill(mock.Anything, mock.Anything).Return(int64(0), apperrors.ErrEventNotFound)

		ctrl.CreateBill(rec, req)

		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("usecase other error", func(t *testing.T) {
		uc := billusecasemock.NewMockUseCase(t)
		ctrl := &controllerImpl{billUC: uc, logger: testkit.NewTestLogger()}
		req := httptest.NewRequest(http.MethodPost, "/bills", bytes.NewBufferString(`{"event_id":1}`))
		req = testkit.WithUserSession(req, vo.UserID(1))
		rec := httptest.NewRecorder()

		uc.EXPECT().CreateBill(mock.Anything, mock.Anything).Return(int64(0), errControllerTest)

		ctrl.CreateBill(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestBillController_FetchEventBills(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		uc := billusecasemock.NewMockUseCase(t)
		ctrl := &controllerImpl{billUC: uc, logger: testkit.NewTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/bills?event_id=5", nil)
		req = testkit.WithUserSession(req, vo.UserID(2))
		rec := httptest.NewRecorder()

		uc.EXPECT().
			FetchEventBills(mock.Anything, int64(5)).
			Return([]domainbill.Bill{
				{
					ID:           1,
					EventID:      5,
					Name:         "Dinner",
					CreatedBy:    vo.UserID(2),
					TotalAmount:  1000,
					Participants: []domainbill.Participant{{MemberID: 10, Amount: 500}},
				},
			}, nil)

		ctrl.FetchEventBills(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		var resp []Bill
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Len(t, resp, 1)
		require.Equal(t, int64(1), resp[0].ID)
		require.Equal(t, int64(10), resp[0].Participants[0].MemberID)
	})

	t.Run("missing session", func(t *testing.T) {
		ctrl := &controllerImpl{billUC: billusecasemock.NewMockUseCase(t), logger: testkit.NewTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/bills?event_id=5", nil)
		rec := httptest.NewRecorder()

		ctrl.FetchEventBills(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("missing event id", func(t *testing.T) {
		ctrl := &controllerImpl{billUC: billusecasemock.NewMockUseCase(t), logger: testkit.NewTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/bills", nil)
		req = testkit.WithUserSession(req, vo.UserID(1))
		rec := httptest.NewRecorder()

		ctrl.FetchEventBills(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid event id", func(t *testing.T) {
		ctrl := &controllerImpl{billUC: billusecasemock.NewMockUseCase(t), logger: testkit.NewTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/bills?event_id=abc", nil)
		req = testkit.WithUserSession(req, vo.UserID(1))
		rec := httptest.NewRecorder()

		ctrl.FetchEventBills(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("usecase event not found", func(t *testing.T) {
		uc := billusecasemock.NewMockUseCase(t)
		ctrl := &controllerImpl{billUC: uc, logger: testkit.NewTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/bills?event_id=5", nil)
		req = testkit.WithUserSession(req, vo.UserID(1))
		rec := httptest.NewRecorder()

		uc.EXPECT().FetchEventBills(mock.Anything, int64(5)).Return(nil, apperrors.ErrEventNotFound)

		ctrl.FetchEventBills(rec, req)

		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		uc := billusecasemock.NewMockUseCase(t)
		ctrl := &controllerImpl{billUC: uc, logger: testkit.NewTestLogger()}
		req := httptest.NewRequest(http.MethodGet, "/bills?event_id=5", nil)
		req = testkit.WithUserSession(req, vo.UserID(1))
		rec := httptest.NewRecorder()

		uc.EXPECT().FetchEventBills(mock.Anything, int64(5)).Return(nil, errControllerTest)

		ctrl.FetchEventBills(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
