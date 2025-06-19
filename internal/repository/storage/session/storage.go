package session

import (
	"context"
	stderrors "errors"

	"billsplitter-monolith/internal/domain/auth"
	"billsplitter-monolith/internal/errors"
	"gorm.io/gorm"
)

type Storage struct {
	db    *gorm.DB
	cache auth.SessionCache
}

func NewStorage(db *gorm.DB, cache auth.SessionCache) auth.SessionStorage {
	return &Storage{
		db:    db,
		cache: cache,
	}
}

func (s *Storage) Create(ctx context.Context, session *auth.Session) error {
	e := fromDomain(session)

	if err := s.db.WithContext(ctx).Create(e).Error; err != nil {
		return errors.ErrSessionStorageFunc(err, "create")
	}

	// TODO: set cache

	return nil
}

func (s *Storage) Get(ctx context.Context, id string) (*auth.Session, error) {
	e := &sessionEntity{}

	// TODO: get cache

	err := s.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS null AND expire_at > now()::timestamptz", id).
		First(&e).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.ErrSessionStorageFunc(err, "get")
	}

	return toDomain(e), nil
}
