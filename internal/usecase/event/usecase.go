package event

import (
	"context"
	"fmt"

	"billsplitter-monolith/internal/domain/event"
	"billsplitter-monolith/internal/domain/user"
	vo "billsplitter-monolith/internal/domain/valueobject"
	"billsplitter-monolith/internal/errors"
	"billsplitter-monolith/internal/utils"
)

type UseCase interface {
	CreateMeet(ctx context.Context, rq CreateMeetRq) (int64, error)
}

type UseCaseImpl struct {
	eventSvc event.Service
	userSvc  user.Service
}

func New(
	eventSvc event.Service,
	userSvc user.Service,
) *UseCaseImpl {
	return &UseCaseImpl{
		eventSvc: eventSvc,
		userSvc:  userSvc,
	}
}

func (uc *UseCaseImpl) CreateMeet(ctx context.Context, rq CreateMeetRq) (int64, error) {
	creator, err := uc.userSvc.GetByID(ctx, rq.CreatedByUserID)
	if err != nil {
		return 0, err
	}

	if creator == nil {
		return 0, fmt.Errorf("get event creator error: %w", errors.ErrUserNotFound)
	}

	meetID, err := uc.eventSvc.Create(ctx, event.CreateEventRq{
		Name:      rq.EventName,
		EventDate: rq.Date,
		Creator: event.Member{
			// сразу привязываем userID для создателя мита
			UserID: utils.Ptr(vo.UserID(creator.ID)),
			Name:   creator.Username,
		},
		Members: rq.Members,
	})

	return meetID, nil
}
