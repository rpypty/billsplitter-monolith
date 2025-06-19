package impl

import (
	"context"
	"log/slog"
	"time"

	"billsplitter-monolith/internal/domain/auth"
	"billsplitter-monolith/internal/errors"
	"billsplitter-monolith/internal/utils"
)

const (
	sessionLiveTime = time.Second * 10
)

type Service struct {
	userStorage    auth.UserStorage
	sessionStorage auth.SessionStorage

	logger *slog.Logger
}

func New(
	userStorage auth.UserStorage,
	sessionStorage auth.SessionStorage,
	logger *slog.Logger,
) *Service {
	return &Service{
		userStorage:    userStorage,
		sessionStorage: sessionStorage,
		logger:         logger,
	}
}

func (s *Service) GetUserBySessionID(ctx context.Context, sessionID string) (*auth.User, error) {
	session, err := s.sessionStorage.Get(ctx, sessionID)
	if err != nil {
		return nil, errors.ErrUserServiceFunc(err, "GetUserBySessionID")
	}

	if session == nil {
		return nil, errors.ErrUserServiceFunc(errors.ErrSessionNotFound, "GetUserBySessionID")
	}

	return s.userStorage.GetByID(ctx, session.UserID)
}

func (s *Service) CreateSession(ctx context.Context, session *auth.Session) (string, error) {
	mth := "CreateSession"

	if session.ID == "" {
		session.ID = utils.NewUUIDv7()
	}

	if session.ExpireAt == nil {
		session.ExpireAt = utils.Ptr(time.Now().Add(sessionLiveTime))
	}

	err := s.sessionStorage.Create(ctx, session)
	if err != nil {
		return "", errors.ErrUserServiceFunc(err, mth)
	}

	return session.ID, nil
}

func (s *Service) CreateOrGetUserByTgID(ctx context.Context, tgID int64, data *auth.User) (*auth.User, error) {
	mth := "CreateOrGetUserByTgID"

	existing, err := s.userStorage.GetByTelegramID(ctx, tgID)
	if err != nil {
		return nil, errors.ErrUserServiceFunc(err, mth)
	}
	if existing != nil {
		return existing, nil
	}

	// validate before create
	if data.ID == "" {
		data.ID = utils.NewUUIDv7()
	}

	if data.Extra.TelegramID == 0 {
		data.Extra.TelegramID = tgID
	}

	err = s.userStorage.Create(ctx, data)
	if err != nil {
		return nil, errors.ErrUserServiceFunc(err, mth)
	}

	u, err := s.userStorage.GetByID(ctx, data.ID)
	if err != nil {
		return nil, errors.ErrUserServiceFunc(err, mth)
	}

	return u, nil
}

func (s *Service) l() *slog.Logger {
	return s.logger.WithGroup("auth-service")
}
