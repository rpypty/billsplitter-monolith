package user

import (
	"context"
	stderrors "errors"

	"gorm.io/gorm"

	"billsplitter-monolith/internal/domain/auth"
	"billsplitter-monolith/internal/errors"
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
	user := &userEntity{}

	if err := s.db.WithContext(ctx).Where("(extra->>'telegram_id')::int = ?", telegramID).First(&user).Error; err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.ErrUserStorageFunc(err, "GetByTelegramID")
	}

	return toDomain(user), nil
}

func (s *Storage) GetByID(ctx context.Context, id string) (*auth.User, error) {
	user := &userEntity{}

	if err := s.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.ErrUserStorageFunc(err, "GetByID")
	}

	return toDomain(user), nil
}

func (s *Storage) Create(ctx context.Context, user *auth.User) error {
	e := fromDomain(user)
	return s.db.WithContext(ctx).Create(e).Error
}
