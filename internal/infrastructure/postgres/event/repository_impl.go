package event

import (
	"context"
	"fmt"

	domain "billsplitter-monolith/internal/domain/event"
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
