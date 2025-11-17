package sdk

import (
	"billsplitter-monolith/internal/transport/http/auth"
	"billsplitter-monolith/internal/transport/http/bill"
	"billsplitter-monolith/internal/transport/http/meet"
	"billsplitter-monolith/internal/transport/http/user"
	hu "billsplitter-monolith/internal/utils/http"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	resty "github.com/go-resty/resty/v2"
)

type TestSDK interface {
	// auth
	Auth(ctx context.Context, rq auth.LoginTelegramReq) (*auth.LoginTelegramRes, error)
	Me(ctx context.Context) (*auth.MeRes, error)

	// meet
	CreateMeet(ctx context.Context, rq meet.CreateEventRq) (*hu.ResponseID, error)
	FetchUserMeets(ctx context.Context) (*[]meet.Event, error)
	GetMeetDetailsByID(ctx context.Context, meetID int64) (*meet.Event, error)
	GetMeetSummary(ctx context.Context, meetID int64) (*meet.EventSummary, error)

	// bill
	CreateBill(ctx context.Context, rq bill.CreateBillRq) (*hu.ResponseID, error)
	FetchEventBills(ctx context.Context, eventID int64) (*[]bill.Bill, error)

	// user
	UpdateUserProfile(ctx context.Context, rq user.UpdateUserProfileReq) (*hu.ResponseOK, error)
	CreatePaymentMethod(ctx context.Context, rq user.CreatePaymentMethodReq) (*user.PaymentMethodResponse, error)
	GetPaymentMethods(ctx context.Context) (*user.PaymentMethodsResponse, error)
	UpdatePaymentMethod(ctx context.Context, methodID int64, rq user.UpdatePaymentMethodReq) (*hu.ResponseOK, error)
	DeletePaymentMethod(ctx context.Context, methodID int64) (*hu.ResponseOK, error)
}

func New(baseUrl string) *SDK {
	restyClient := resty.New().SetBaseURL(baseUrl)

	return &SDK{
		client: restyClient,
		mu:     &sync.Mutex{},
	}
}

type SDK struct {
	sessionID string
	userID    int64

	mu     *sync.Mutex
	client *resty.Client
}

func (s *SDK) Auth(ctx context.Context, rq auth.LoginTelegramReq) (*auth.LoginTelegramRes, error) {
	rawResp, err := s.client.
		R().
		SetContext(ctx).
		SetBody(&rq).
		Post("/auth/login/telegram")

	if err != nil {
		return nil, err
	}

	var resp auth.LoginTelegramRes

	err = json.Unmarshal(rawResp.Body(), &resp)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessionID = resp.SessionID
	s.userID = int64(resp.UserInfo.ID)

	return &resp, nil
}

func (s *SDK) Me(ctx context.Context) (*auth.MeRes, error) {
	rawResp, err := s.newRequest(ctx).Get("/auth/me")

	if err != nil {
		return nil, err
	}

	var resp auth.MeRes

	err = json.Unmarshal(rawResp.Body(), &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *SDK) CreateMeet(ctx context.Context, rq meet.CreateEventRq) (*hu.ResponseID, error) {
	rawResp, err := s.newRequest(ctx).
		SetBody(&rq).
		Post("/meets")
	if err != nil {
		return nil, err
	}

	var resp hu.ResponseID
	if err := json.Unmarshal(rawResp.Body(), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *SDK) FetchUserMeets(ctx context.Context) (*[]meet.Event, error) {
	rawResp, err := s.newRequest(ctx).Get("/meets")
	if err != nil {
		return nil, err
	}

	var resp []meet.Event
	if err := json.Unmarshal(rawResp.Body(), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *SDK) GetMeetDetailsByID(ctx context.Context, meetID int64) (*meet.Event, error) {
	rawResp, err := s.newRequest(ctx).Get(fmt.Sprintf("/meets/%d", meetID))
	if err != nil {
		return nil, err
	}

	var resp meet.Event
	if err := json.Unmarshal(rawResp.Body(), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *SDK) GetMeetSummary(ctx context.Context, meetID int64) (*meet.EventSummary, error) {
	rawResp, err := s.newRequest(ctx).Get(fmt.Sprintf("/meets/%d/summary", meetID))
	if err != nil {
		return nil, err
	}

	var resp meet.EventSummary
	if err := json.Unmarshal(rawResp.Body(), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *SDK) CreateBill(ctx context.Context, rq bill.CreateBillRq) (*hu.ResponseID, error) {
	rawResp, err := s.newRequest(ctx).
		SetBody(&rq).
		Post("/bills")
	if err != nil {
		return nil, err
	}

	var resp hu.ResponseID
	if err := json.Unmarshal(rawResp.Body(), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *SDK) FetchEventBills(ctx context.Context, eventID int64) (*[]bill.Bill, error) {
	rawResp, err := s.newRequest(ctx).
		SetQueryParam("event_id", fmt.Sprintf("%d", eventID)).
		Get("/bills")
	if err != nil {
		return nil, err
	}

	var resp []bill.Bill
	if err := json.Unmarshal(rawResp.Body(), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *SDK) UpdateUserProfile(ctx context.Context, rq user.UpdateUserProfileReq) (*hu.ResponseOK, error) {
	rawResp, err := s.newRequest(ctx).
		SetBody(&rq).
		Put(fmt.Sprintf("/user/%d/profile", s.getUserID()))
	if err != nil {
		return nil, err
	}

	var resp hu.ResponseOK
	if err := json.Unmarshal(rawResp.Body(), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *SDK) CreatePaymentMethod(ctx context.Context, rq user.CreatePaymentMethodReq) (*user.PaymentMethodResponse, error) {
	rawResp, err := s.newRequest(ctx).
		SetBody(&rq).
		Post("/payment_methods")
	if err != nil {
		return nil, err
	}

	var resp user.PaymentMethodResponse
	if err := json.Unmarshal(rawResp.Body(), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *SDK) GetPaymentMethods(ctx context.Context) (*user.PaymentMethodsResponse, error) {
	rawResp, err := s.newRequest(ctx).Get("/payment_methods")
	if err != nil {
		return nil, err
	}

	var resp user.PaymentMethodsResponse
	if err := json.Unmarshal(rawResp.Body(), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *SDK) UpdatePaymentMethod(ctx context.Context, methodID int64, rq user.UpdatePaymentMethodReq) (*hu.ResponseOK, error) {
	rawResp, err := s.newRequest(ctx).
		SetBody(&rq).
		Put(fmt.Sprintf("/payment_methods/%d", methodID))
	if err != nil {
		return nil, err
	}

	var resp hu.ResponseOK
	if err := json.Unmarshal(rawResp.Body(), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *SDK) DeletePaymentMethod(ctx context.Context, methodID int64) (*hu.ResponseOK, error) {
	rawResp, err := s.newRequest(ctx).
		Delete(fmt.Sprintf("/payment_methods/%d", methodID))
	if err != nil {
		return nil, err
	}

	var resp hu.ResponseOK
	if err := json.Unmarshal(rawResp.Body(), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *SDK) newRequest(ctx context.Context) *resty.Request {
	req := s.client.R().SetContext(ctx)

	s.mu.Lock()
	sessionID := s.sessionID
	s.mu.Unlock()

	if sessionID != "" {
		req = req.SetHeader("X-Session-ID", sessionID)
	}

	return req
}

func (s *SDK) getUserID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.userID
}
