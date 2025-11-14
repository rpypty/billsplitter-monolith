package user

import (
	"context"
	stderrors "errors"
	"fmt"

	domain "billsplitter-monolith/internal/domain/user"
	vo "billsplitter-monolith/internal/domain/valueobject"
	"gorm.io/gorm"
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

func (s *RepositoryImpl) Fetch(ctx context.Context, filter domain.FetchFilter) ([]domain.User, error) {
	errWrap := getErrWrapper("Fetch")

	var users []userEntity

	err := s.db.WithContext(ctx).Where("id IN ?", filter.UserIDs).Find(&users).Error
	if err != nil {
		return nil, errWrap(err)
	}

	domainUsers := make([]domain.User, 0, len(users))
	for _, user := range users {
		domainUsers = append(domainUsers, *toDomain(&user))
	}

	return domainUsers, nil
}

func (s *RepositoryImpl) GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	errWrap := getErrWrapper("GetByTelegramID")

	user := &userEntity{}

	err := s.db.
		WithContext(ctx).
		Where("(extra->>'telegram_id')::bigint = ?", telegramID).
		First(&user).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errWrap(err)
	}

	return toDomain(user), nil
}

func (s *RepositoryImpl) GetByID(ctx context.Context, id vo.UserID) (*domain.User, error) {
	errWrap := getErrWrapper("GetByID")

	user := &userEntity{}

	err := s.db.WithContext(ctx).First(&user, "id = ?", int64(id)).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, errWrap(err)
	}

	return toDomain(user), nil
}

func (s *RepositoryImpl) Update(ctx context.Context, id vo.UserID, user domain.User) error {
	errWrap := getErrWrapper("Update")

	err := s.db.
		WithContext(ctx).
		Model(&userEntity{}).
		Where("id = ?", int64(id)).
		Updates(map[string]any{
			"username":   user.Username,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		}).Error
	if err != nil {
		return errWrap(err)
	}

	return nil
}

func (s *RepositoryImpl) Create(ctx context.Context, user domain.User) (*domain.User, error) {
	errWrap := getErrWrapper("Create")

	e := fromDomain(&user)

	err := s.db.WithContext(ctx).Create(e).Error
	if err != nil {
		return nil, errWrap(err)
	}

	return toDomain(e), nil
}

func getErrWrapper(method string) func(error) error {
	return func(err error) error {
		return fmt.Errorf("UserRepositoryIpml->%s: %w", method, err)
	}
}
