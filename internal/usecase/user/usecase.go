package user

import (
	"context"
	stderrors "errors"

	"billsplitter-monolith/internal/domain/session"
	"billsplitter-monolith/internal/domain/user"
	"billsplitter-monolith/internal/errors"
)

// UseCase - верхнеуровневый сервис, который оркестрирует domain сервисами
type UseCase interface {
	GetBySessionID(ctx context.Context, sessionID string) (*user.User, error)

	// GetOrCreate - атомарно проверяет есть ли пользователь, если нет - создает нового
	GetByTgIDOrCreate(ctx context.Context, user user.User) (*user.User, error)
}

type UseCaseImpl struct {
	sessionSvc session.Service
	userSvc    user.Service
}

func New(
	userSvc user.Service,
	sessionSvc session.Service,
) *UseCaseImpl {
	return &UseCaseImpl{
		userSvc:    userSvc,
		sessionSvc: sessionSvc,
	}
}

func (uc *UseCaseImpl) GetByTgIDOrCreate(ctx context.Context, user user.User) (*user.User, error) {
	// TODO: надо добавить UseCase транзакцию бдшки для атомарности

	userInfo, err := uc.userSvc.GetByTelegramID(ctx, user.Extra.TelegramID)
	if err == nil {
		return userInfo, nil
	}

	if !stderrors.Is(err, errors.ErrUserNotFound) {
		// some internal error
		return nil, err
	}

	// user not found -> create

	createUserInfo, err := uc.userSvc.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return createUserInfo, nil
}

// GetBySessionID - возвращает юзера по sessionID
func (uc *UseCaseImpl) GetBySessionID(ctx context.Context, sessionID string) (*user.User, error) {
	sessionInfo, err := uc.sessionSvc.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	userInfo, err := uc.userSvc.GetByID(ctx, sessionInfo.UserID)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}
