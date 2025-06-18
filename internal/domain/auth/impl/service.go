package impl

import (
	"context"
	"log/slog"

	"billsplitter-monolith/internal/domain/auth"
	"billsplitter-monolith/internal/errors"
	"billsplitter-monolith/internal/utils"
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
		return nil, err
	}

	if session == nil {
		return nil, errors.ErrSessionNotFound
	}

	return s.userStorage.GetByID(ctx, session.UserID)
}

func (s *Service) CreateSession(ctx context.Context, session *auth.Session) (string, error) {
	if session.ID == "" {
		session.ID = utils.NewUUIDv7()
	}

	err := s.sessionStorage.Create(ctx, session.UserID)
	if err != nil {
		return "", err
	}

	return session.ID, nil
}

func (s *Service) CreateOrGetUserByTgID(ctx context.Context, tgID int64, data *auth.User) (*auth.User, error) {
	existing, err := s.userStorage.GetByTelegramID(ctx, tgID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	// validate before create
	if data.ID == "" {
		data.ID = utils.NewUUIDv7()
	}

	err = s.userStorage.Create(ctx, data)
	if err != nil {
		return nil, err
	}

	return s.userStorage.GetByID(ctx, data.ID)
}

func (s *Service) l() *slog.Logger {
	return s.logger.WithGroup("auth-service")
}
