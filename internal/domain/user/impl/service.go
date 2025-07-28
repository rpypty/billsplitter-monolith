package impl

import (
	"context"
	stderrors "errors"
	"fmt"

	domain "billsplitter-monolith/internal/domain/user"
	"billsplitter-monolith/internal/errors"
)

var _ domain.Service = (*ServiceImpl)(nil)

type ServiceImpl struct {
	repo domain.Repository
}

func New(repo domain.Repository) *ServiceImpl {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) Fetch(ctx context.Context, filter domain.FetchFilter) ([]domain.User, error) {
	errWrap := getErrWrapperFunc("Fetch")

	users, err := s.repo.Fetch(ctx, filter)
	if err != nil {
		return nil, errWrap(err)
	}

	return users, nil
}

func (s *ServiceImpl) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	errWrap := getErrWrapperFunc("GetById")

	// TODO: add cache layer

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errWrap(err)
	}

	if user == nil {
		return nil, errWrap(stderrors.New("user not found"))
	}

	return user, nil
}

func (s *ServiceImpl) GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	errWrap := getErrWrapperFunc("GetByTelegramID")

	// TODO: add cache layer by telegramID

	user, err := s.repo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, errWrap(err)
	}

	if user == nil {
		return nil, errWrap(errors.ErrUserNotFound)
	}

	return user, nil
}

func (s *ServiceImpl) Create(ctx context.Context, user domain.User) (*domain.User, error) {
	errWrap := getErrWrapperFunc("Create")

	err := validateUser(user)
	if err != nil {
		return nil, errWrap(err)
	}

	userInfo, err := s.repo.GetByTelegramID(ctx, user.Extra.TelegramID)
	if err != nil {
		return nil, errWrap(err)
	}

	if userInfo != nil {
		return nil, errWrap(stderrors.New("this telegramID already exists"))
	}

	created, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, errWrap(err)
	}

	// TODO: add cache set

	return created, nil
}

func (s *ServiceImpl) Update(ctx context.Context, id int64, user domain.User) error {
	errWrap := getErrWrapperFunc("Update")

	err := validateUser(user)
	if err != nil {
		return errWrap(err)
	}

	err = s.repo.Update(ctx, id, user)
	if err != nil {
		return errWrap(err)
	}

	// TODO: add cache set

	return nil
}

func validateUser(user domain.User) error {
	if user.Username == "" {
		return errors.ErrValidationFunc("user name is required")
	}

	return nil
}

func getErrWrapperFunc(mth string) func(error) error {
	return func(err error) error {
		return fmt.Errorf("domain/user/impl/ServiceImpl.%s: %w", mth, err)
	}
}
