package bill

import (
	"context"
	"fmt"

	domain "billsplitter-monolith/internal/domain/bill"

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

func (r *RepositoryImpl) Create(ctx context.Context, b domain.Bill) (int64, error) {
	errWrap := getErrWrapper("Create")

	entity := billFromDomain(&b)

	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return 0, errWrap(err)
	}

	return entity.ID, nil
}

func (r *RepositoryImpl) FetchByEventID(ctx context.Context, eventID int64) ([]domain.Bill, error) {
	errWrap := getErrWrapper("FetchByEventID")

	entities := make([]billEntity, 0)

	err := r.db.
		WithContext(ctx).
		Preload("Participants").
		Where("event_id = ? AND deleted_at IS NULL", eventID).
		Order("created_at DESC").
		Find(&entities).Error
	if err != nil {
		return nil, errWrap(err)
	}

	bills := make([]domain.Bill, 0, len(entities))

	for _, entity := range entities {
		if d := billToDomain(&entity); d != nil {
			bills = append(bills, *d)
		}
	}

	return bills, nil
}

func getErrWrapper(method string) func(error) error {
	return func(err error) error {
		return fmt.Errorf("BillRepositoryImpl->%s: %w", method, err)
	}
}
