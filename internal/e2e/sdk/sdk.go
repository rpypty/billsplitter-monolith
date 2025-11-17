package sdk

import (
	"billsplitter-monolith/internal/transport/http/auth"
	"billsplitter-monolith/internal/transport/http/meet"
	hu "billsplitter-monolith/internal/utils/http"
	"context"
	"encoding/json"
	"sync"

	resty "github.com/go-resty/resty/v2"
)

type TestSDK interface {
	// auth
	Auth(ctx context.Context, rq auth.LoginTelegramReq) (*auth.LoginTelegramRes, error)
	Me(ctx context.Context) (*auth.MeRes, error)

	// user
	CreateMeet(ctx context.Context, rq meet.CreateEventRq) (*hu.ResponseID, error)
	FetchUserMeets(ctx context.Context) (*[]meet.Event, error)
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
	req := s.client.
		R().
		SetContext(ctx)

	if s.sessionID != "" {
		req = req.SetHeader("X-Session-ID", s.sessionID)
	}

	rawResp, err := req.Get("/auth/me")

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
