package impl

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"billsplitter-monolith/internal/domain/session"
	vo "billsplitter-monolith/internal/domain/valueobject"
	"billsplitter-monolith/internal/errors"
	"billsplitter-monolith/internal/utils"
)

const (
	sessionLiveTime = time.Minute * 60
)

var _ session.Service = (*ServiceImpl)(nil)

type ServiceImpl struct {
	repo   session.Repository
	logger *slog.Logger
}

func New(
	repo session.Repository,
	logger *slog.Logger,
) *ServiceImpl {
	return &ServiceImpl{
		repo:   repo,
		logger: logger,
	}
}

func (s *ServiceImpl) GetByID(ctx context.Context, sessionID string) (*session.Session, error) {
	errWrap := getErrWrapperFunc("GetByID")

	obj, err := s.repo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, errWrap(err)
	}

	if obj == nil {
		return nil, errWrap(errors.ErrSessionNotFound)
	}

	if time.Now().After(*obj.ExpireAt) {
		return nil, errWrap(errors.ErrSessionExpired)
	}

	return obj, nil
}

func (s *ServiceImpl) Create(ctx context.Context, userID vo.UserID) (string, error) {
	errWrap := getErrWrapperFunc("Create")

	if userID == 0 {
		return "", errWrap(errors.ErrValidationFunc("userID is required"))
	}

	obj := session.Session{
		UserID:   userID,
		ExpireAt: utils.Ptr(time.Now().Add(sessionLiveTime)),
		ID:       utils.NewUUIDv7(),
	}

	created, err := s.repo.Create(ctx, obj)
	if err != nil {
		return "", errWrap(err)
	}

	return created.ID, nil
}

func (s *ServiceImpl) l() *slog.Logger {
	return s.logger.WithGroup("auth-service")
}

func getErrWrapperFunc(mth string) func(error) error {
	return func(err error) error {
		return fmt.Errorf("domain/session/impl/ServiceImpl.%s: %w", mth, err)
	}
}
