package session

import (
	"context"
	"errors"
	"time"

	"billsplitter-monolith/internal/domain/auth"
	"gorm.io/gorm"
)

const (
	sessionLiveTime = time.Hour
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

func (s *Storage) Create(ctx context.Context, userID string) error {
	expire := time.Now().Add(sessionLiveTime)

	session := &auth.Session{
		UserID:   userID,
		ExpireAt: &expire,
	}

	if err := s.db.WithContext(ctx).Create(session).Error; err != nil {
		return err
	}

	// TODO: set cache

	return nil
}

func (s *Storage) Get(ctx context.Context, id string) (*auth.Session, error) {
	var session auth.Session

	// TODO: get cache

	err := s.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS null", id).
		First(&session).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &session, nil
}
