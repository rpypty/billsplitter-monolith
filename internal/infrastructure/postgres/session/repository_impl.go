package session

import (
	"context"
	stderrors "errors"
	"fmt"

	"gorm.io/gorm"

	domain "billsplitter-monolith/internal/domain/session"
)

var _ domain.Repository = (*RepositoryImpl)(nil)

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.Repository {
	return &RepositoryImpl{
		db: db,
	}
}

func (s *RepositoryImpl) Create(ctx context.Context, session domain.Session) (*domain.Session, error) {
	errWrap := getErrWrapper("Create")

	entity := fromDomain(&session)

	err := s.db.WithContext(ctx).Create(entity).Error
	if err != nil {
		return nil, errWrap(err)
	}

	return toDomain(entity), nil
}

func (s *RepositoryImpl) GetByID(ctx context.Context, id string) (*domain.Session, error) {
	errWrap := getErrWrapper("GetByID")

	e := &sessionEntity{}

	err := s.db.
		WithContext(ctx).
		Where("id = ? AND deleted_at IS null AND expire_at > now()::timestamptz", id).
		First(&e).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errWrap(err)
	}

	return toDomain(e), nil
}

func getErrWrapper(method string) func(error) error {
	return func(err error) error {
		return fmt.Errorf("SessionRepositoryIpml->%s: %w", method, err)
	}
}
