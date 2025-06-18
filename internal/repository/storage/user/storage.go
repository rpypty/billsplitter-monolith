package user

import (
	"context"
	"errors"

	"billsplitter-monolith/internal/domain/auth"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) auth.UserStorage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) GetByTelegramID(ctx context.Context, telegramID int64) (*auth.User, error) {
	user := auth.User{}

	if err := s.db.WithContext(ctx).Where("telegram_id = ?", telegramID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (s *Storage) GetByID(ctx context.Context, id string) (*auth.User, error) {
	user := auth.User{}

	if err := s.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (s *Storage) Create(ctx context.Context, user *auth.User) error {
	return s.db.WithContext(ctx).Create(user).Error
}
