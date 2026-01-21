package event

import (
	"context"
	"errors"
	"fmt"

	domain "billsplitter-monolith/internal/domain/event"
	vo "billsplitter-monolith/internal/domain/valueobject"
	apperrors "billsplitter-monolith/internal/errors"

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

func (r *RepositoryImpl) Create(ctx context.Context, ev domain.Event) (int64, error) {
	errWrap := getErrWrapper("Create")

	evEntity := eventFromDomain(&ev)

	err := r.db.WithContext(ctx).Create(&evEntity).Error
	if err != nil {
		return 0, errWrap(err)
	}

	return evEntity.ID, nil
}

func (r *RepositoryImpl) AddMembers(ctx context.Context, eventID int64, members []domain.Member) error {
	errWrap := getErrWrapper("AddMembers")

	entityMembers := make([]*memberEntity, 0, len(members))

	for _, member := range members {
		entityMembers = append(entityMembers, memberFromDomain(eventID, &member))
	}

	err := r.db.WithContext(ctx).Create(entityMembers).Error
	if err != nil {
		return errWrap(err)
	}

	return nil
}

func getErrWrapper(method string) func(error) error {
	return func(err error) error {
		return fmt.Errorf("EventRepositoryIpml->%s: %w", method, err)
	}
}

func (r *RepositoryImpl) FetchByUserID(ctx context.Context, userID vo.UserID) ([]domain.Event, error) {
	entityEvents := make([]eventEntity, 0)

	err := r.db.
		WithContext(ctx).
		Preload("Members").
		Order("created_at DESC").
		Find(&entityEvents, "created_by_user_id = ?", int64(userID)).
		Error
	if err != nil {
		return nil, err
	}

	domainEvents := make([]domain.Event, 0, len(entityEvents))

	for _, ev := range entityEvents {
		domainEvents = append(domainEvents, *eventToDomain(&ev))
	}

	return domainEvents, nil
}

func (r *RepositoryImpl) GetByID(ctx context.Context, eventID int64) (*domain.Event, error) {
	entityEvent := &eventEntity{}

	err := r.db.
		WithContext(ctx).
		Preload("Members").
		Find(&entityEvent, "id = ?", eventID).
		Error
	if err != nil {
		return nil, err
	}

	return eventToDomain(entityEvent), nil
}

func (r *RepositoryImpl) GetByPublicUUID(ctx context.Context, publicUUID string) (*domain.Event, error) {
	entityEvent := &eventEntity{}

	err := r.db.
		WithContext(ctx).
		Preload("Members").
		Find(&entityEvent, "public_uuid = ?", publicUUID).
		Error
	if err != nil {
		return nil, err
	}

	return eventToDomain(entityEvent), nil
}

func (r *RepositoryImpl) AssignMemberUser(ctx context.Context, eventID int64, memberID int64, userID vo.UserID) error {
	errWrap := getErrWrapper("AssignMemberUser")

	result := r.db.WithContext(ctx).
		Model(&memberEntity{}).
		Where("id = ? AND event_id = ? AND user_id IS NULL", memberID, eventID).
		Update("user_id", int64(userID))
	if result.Error != nil {
		return errWrap(result.Error)
	}

	if result.RowsAffected == 0 {
		var member memberEntity
		err := r.db.WithContext(ctx).
			Select("id", "user_id").
			Where("id = ? AND event_id = ?", memberID, eventID).
			Take(&member).
			Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperrors.ErrValidationFunc("member not found in event")
			}
			return errWrap(err)
		}

		if member.UserID != nil {
			return apperrors.ErrValidationFunc("member already assigned")
		}
	}

	return nil
}
