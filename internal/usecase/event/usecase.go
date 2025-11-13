package event

import (
	"context"
	"fmt"

	"billsplitter-monolith/internal/domain/event"
	"billsplitter-monolith/internal/domain/user"
	vo "billsplitter-monolith/internal/domain/valueobject"
	"billsplitter-monolith/internal/errors"
	"billsplitter-monolith/internal/transport/http/middleware"
)

type UseCase interface {
	CreateMeet(ctx context.Context, rq CreateMeetRq) (int64, error)
	FetchUserMeets(ctx context.Context, userID vo.UserID) ([]event.Event, error)
	GetMeetByID(ctx context.Context, meetID int64) (*event.Event, error)
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
	creatorUser, err := uc.userSvc.GetByID(ctx, rq.CreatedByUserID)
	if err != nil {
		return 0, err
	}

	if creatorUser == nil {
		return 0, fmt.Errorf("get event creator error: %w", errors.ErrUserNotFound)
	}

	meetID, err := uc.eventSvc.Create(ctx, event.CreateEventRq{
		Name:            rq.EventName,
		EventDate:       rq.Date,
		CreatorUserID:   creatorUser.ID,
		CreatorUsername: creatorUser.Username,
		Members:         rq.Members,
	})
	if err != nil {
		return 0, err
	}

	return meetID, nil
}

func (uc *UseCaseImpl) FetchUserMeets(ctx context.Context, userID vo.UserID) ([]event.Event, error) {
	events, err := uc.eventSvc.FetchByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (uc *UseCaseImpl) GetMeetByID(ctx context.Context, meetID int64) (*event.Event, error) {
	event, err := uc.eventSvc.GetByID(ctx, meetID)
	if err != nil {
		return nil, err
	}

	if event == nil || event.ID == 0 {
		return nil, nil
	}

	userID, err := middleware.ExtractUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := uc.validateEvent(event, userID); err != nil {
		return nil, err
	}

	return event, nil
}

func (uc *UseCaseImpl) validateEvent(ev *event.Event, sessionUserID vo.UserID) error {
	userMemberIndex := make(map[vo.UserID]struct{}, len(ev.Members))
	for _, m := range ev.Members {
		if m.UserID == nil {
			continue
		}
		userMemberIndex[*m.UserID] = struct{}{}
	}

	// проверяем есть ли инициатор запроса в этом ивенте
	if _, ok := userMemberIndex[sessionUserID]; !ok {
		return errors.ErrForbiden
	}

	return nil
}
