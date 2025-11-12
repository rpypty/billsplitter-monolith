package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	domainpayment "billsplitter-monolith/internal/domain/payment_method"
	domainuser "billsplitter-monolith/internal/domain/user"
	vo "billsplitter-monolith/internal/domain/valueobject"
	paymentmethodmock "billsplitter-monolith/internal/mocks/domain/payment_method"
	usermock "billsplitter-monolith/internal/mocks/domain/user"
	"billsplitter-monolith/internal/utils/testkit"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserController_UpdateUserProfile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userSvc := usermock.NewMockService(t)
		ctrl := &controllerImpl{
			userSvc: userSvc,
			logger:  testkit.NewTestLogger(),
		}
		body := bytes.NewBufferString(`{"username":"nick","firstName":"Nick","lastName":"Doe"}`)
		req := httptest.NewRequest(http.MethodPut, "/user/1/profile", body)
		req = testkit.WithUserSession(req, vo.UserID(1))
		req = testkit.WithRouteParams(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		userSvc.EXPECT().
			Update(mock.Anything, vo.UserID(1), mock.MatchedBy(func(u domainuser.User) bool {
				return u.Username == "nick" && u.FirstName == "Nick" && u.LastName == "Doe"
			})).
			Return(nil)

		ctrl.UpdateUserProfile(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("forbidden on mismatch", func(t *testing.T) {
		ctrl := &controllerImpl{
			userSvc: usermock.NewMockService(t),
			logger:  testkit.NewTestLogger(),
		}
		req := httptest.NewRequest(http.MethodPut, "/user/2/profile", bytes.NewBufferString(`{}`))
		req = testkit.WithUserSession(req, vo.UserID(1))
		req = testkit.WithRouteParams(req, map[string]string{"id": "2"})
		rec := httptest.NewRecorder()

		ctrl.UpdateUserProfile(rec, req)

		require.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("invalid path id", func(t *testing.T) {
		ctrl := &controllerImpl{
			userSvc: usermock.NewMockService(t),
			logger:  testkit.NewTestLogger(),
		}
		req := httptest.NewRequest(http.MethodPut, "/user/foo/profile", bytes.NewBufferString(`{}`))
		req = testkit.WithUserSession(req, vo.UserID(1))
		req = testkit.WithRouteParams(req, map[string]string{"id": "foo"})
		rec := httptest.NewRecorder()

		ctrl.UpdateUserProfile(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestUserController_GetPaymentMethods(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userSvc := usermock.NewMockService(t)
		paymentSvc := paymentmethodmock.NewMockService(t)
		ctrl := &controllerImpl{
			userSvc:          userSvc,
			paymentMethodSvc: paymentSvc,
			logger:           testkit.NewTestLogger(),
		}
		req := httptest.NewRequest(http.MethodGet, "/user/1/payment_methods", nil)
		req = testkit.WithUserSession(req, vo.UserID(1))
		req = testkit.WithRouteParams(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		methods := []domainpayment.PaymentMethod{
			{ID: 1, UserID: 1, Name: "SBP", Description: "desc", Recipient: "+7999"},
		}

		paymentSvc.EXPECT().
			FetchByUserID(mock.Anything, vo.UserID(1)).
			Return(methods, nil)

		ctrl.GetPaymentMethods(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		var resp PaymentMethodsResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Len(t, resp.PaymentMethods, 1)
		require.Equal(t, "SBP", resp.PaymentMethods[0].Name)
	})

	t.Run("invalid path id", func(t *testing.T) {
		ctrl := &controllerImpl{
			userSvc:          usermock.NewMockService(t),
			paymentMethodSvc: paymentmethodmock.NewMockService(t),
			logger:           testkit.NewTestLogger(),
		}
		req := httptest.NewRequest(http.MethodGet, "/user/foo/payment_methods", nil)
		req = testkit.WithUserSession(req, vo.UserID(1))
		req = testkit.WithRouteParams(req, map[string]string{"id": "foo"})
		rec := httptest.NewRecorder()

		ctrl.GetPaymentMethods(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("forbidden when session mismatch", func(t *testing.T) {
		ctrl := &controllerImpl{
			userSvc:          usermock.NewMockService(t),
			paymentMethodSvc: paymentmethodmock.NewMockService(t),
			logger:           testkit.NewTestLogger(),
		}
		req := httptest.NewRequest(http.MethodGet, "/user/2/payment_methods", nil)
		req = testkit.WithUserSession(req, vo.UserID(1))
		req = testkit.WithRouteParams(req, map[string]string{"id": "2"})
		rec := httptest.NewRecorder()

		ctrl.GetPaymentMethods(rec, req)

		require.Equal(t, http.StatusForbidden, rec.Code)
	})
}

func TestUserController_CreatePaymentMethod(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		paymentSvc := paymentmethodmock.NewMockService(t)
		ctrl := &controllerImpl{
			paymentMethodSvc: paymentSvc,
			logger:           testkit.NewTestLogger(),
		}
		body := bytes.NewBufferString(`{"name":"SBP","description":"desc","recipient":"+7999"}`)
		req := httptest.NewRequest(http.MethodPost, "/payment_methods", body)
		req = testkit.WithUserSession(req, vo.UserID(5))
		rec := httptest.NewRecorder()

		created := domainpayment.PaymentMethod{ID: 10, UserID: 5, Name: "SBP", Description: "desc", Recipient: "+7999"}

		paymentSvc.EXPECT().
			Create(mock.Anything, mock.MatchedBy(func(pm domainpayment.PaymentMethod) bool {
				return pm.UserID == 5 && pm.Name == "SBP"
			})).
			Return(created, nil)

		ctrl.CreatePaymentMethod(rec, req)

		require.Equal(t, http.StatusCreated, rec.Code)
		var resp PaymentMethodResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, int64(10), resp.ID)
	})

	t.Run("missing required fields", func(t *testing.T) {
		ctrl := &controllerImpl{
			paymentMethodSvc: paymentmethodmock.NewMockService(t),
			logger:           testkit.NewTestLogger(),
		}
		req := httptest.NewRequest(http.MethodPost, "/payment_methods", bytes.NewBufferString(`{"description":"desc"}`))
		req = testkit.WithUserSession(req, vo.UserID(5))
		rec := httptest.NewRecorder()

		ctrl.CreatePaymentMethod(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("forbidden", func(t *testing.T) {
		ctrl := &controllerImpl{
			paymentMethodSvc: paymentmethodmock.NewMockService(t),
			logger:           testkit.NewTestLogger(),
		}
		req := httptest.NewRequest(http.MethodPost, "/user/2/payment_methods", bytes.NewBufferString(`{"name":"SBP","recipient":"r"}`))
		req = testkit.WithUserSession(req, vo.UserID(1))
		req = testkit.WithRouteParams(req, map[string]string{"id": "2"})
		rec := httptest.NewRecorder()

		ctrl.CreatePaymentMethod(rec, req)

		require.Equal(t, http.StatusForbidden, rec.Code)
	})
}

func TestUserController_UpdatePaymentMethod(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		paymentSvc := paymentmethodmock.NewMockService(t)
		ctrl := &controllerImpl{
			paymentMethodSvc: paymentSvc,
			logger:           testkit.NewTestLogger(),
		}
		body := bytes.NewBufferString(`{"name":"SBP","description":"desc","recipient":"+7999"}`)
		req := httptest.NewRequest(http.MethodPut, "/payment_methods/1", body)
		req = testkit.WithUserSession(req, vo.UserID(1))
		req = testkit.WithRouteParams(req, map[string]string{"id": "1", "methodId": "9"})
		rec := httptest.NewRecorder()

		paymentSvc.EXPECT().
			Update(mock.Anything, int64(9), mock.MatchedBy(func(pm domainpayment.PaymentMethod) bool {
				return pm.Name == "SBP" && pm.UserID == 1
			})).
			Return(nil)

		ctrl.UpdatePaymentMethod(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("invalid method id", func(t *testing.T) {
		ctrl := &controllerImpl{
			paymentMethodSvc: paymentmethodmock.NewMockService(t),
			logger:           testkit.NewTestLogger(),
		}
		req := httptest.NewRequest(http.MethodPut, "/payment_methods/foo", bytes.NewBufferString(`{}`))
		req = testkit.WithUserSession(req, vo.UserID(1))
		req = testkit.WithRouteParams(req, map[string]string{"methodId": "foo"})
		rec := httptest.NewRecorder()

		ctrl.UpdatePaymentMethod(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestUserController_DeletePaymentMethod(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		paymentSvc := paymentmethodmock.NewMockService(t)
		ctrl := &controllerImpl{
			paymentMethodSvc: paymentSvc,
			logger:           testkit.NewTestLogger(),
		}
		req := httptest.NewRequest(http.MethodDelete, "/payment_methods/9", nil)
		req = testkit.WithUserSession(req, vo.UserID(1))
		req = testkit.WithRouteParams(req, map[string]string{"id": "1", "methodId": "9"})
		rec := httptest.NewRecorder()

		paymentSvc.EXPECT().
			Delete(mock.Anything, int64(9)).
			Return(nil)

		ctrl.DeletePaymentMethod(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("forbidden", func(t *testing.T) {
		ctrl := &controllerImpl{
			paymentMethodSvc: paymentmethodmock.NewMockService(t),
			logger:           testkit.NewTestLogger(),
		}
		req := httptest.NewRequest(http.MethodDelete, "/user/2/payment_methods/9", nil)
		req = testkit.WithUserSession(req, vo.UserID(1))
		req = testkit.WithRouteParams(req, map[string]string{"id": "2", "methodId": "9"})
		rec := httptest.NewRecorder()

		ctrl.DeletePaymentMethod(rec, req)

		require.Equal(t, http.StatusForbidden, rec.Code)
	})
}
