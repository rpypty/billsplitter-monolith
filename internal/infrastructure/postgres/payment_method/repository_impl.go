package payment_method

import (
	"context"
	stderrors "errors"
	"fmt"

	domain "billsplitter-monolith/internal/domain/payment_method"
	vo "billsplitter-monolith/internal/domain/valueobject"
	"gorm.io/gorm"
)

var _ domain.Repository = (*RepositoryImpl)(nil)

type RepositoryImpl struct {
	db *gorm.DB
}

func New(db *gorm.DB) domain.Repository {
	return &RepositoryImpl{
		db: db,
	}
}

func (r *RepositoryImpl) GetByID(ctx context.Context, id int64) (domain.PaymentMethod, error) {
	errWrap := getErrWrapper("GetByID")

	e := &paymentMethodEntity{}

	err := r.db.WithContext(ctx).First(e, "id = ?", id).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return domain.PaymentMethod{}, nil
		}
		return domain.PaymentMethod{}, errWrap(err)
	}

	return *toDomain(e), nil
}

func (r *RepositoryImpl) FetchByUserID(ctx context.Context, userID vo.UserID) ([]domain.PaymentMethod, error) {
	errWrap := getErrWrapper("FetchByUserID")

	var entities []paymentMethodEntity

	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&entities).Error
	if err != nil {
		return nil, errWrap(err)
	}

	methods := make([]domain.PaymentMethod, 0, len(entities))
	for _, entity := range entities {
		methods = append(methods, *toDomain(&entity))
	}

	return methods, nil
}

func (r *RepositoryImpl) Create(ctx context.Context, method domain.PaymentMethod) (domain.PaymentMethod, error) {
	errWrap := getErrWrapper("Create")

	e := fromDomain(&method)

	err := r.db.WithContext(ctx).Create(e).Error
	if err != nil {
		return domain.PaymentMethod{}, errWrap(err)
	}

	return *toDomain(e), nil
}

func (r *RepositoryImpl) Update(ctx context.Context, id int64, method domain.PaymentMethod) error {
	errWrap := getErrWrapper("Update")

	err := r.db.
		WithContext(ctx).
		Model(&paymentMethodEntity{}).
		Where("id = ? AND user_id = ?", id, method.UserID).
		Updates(map[string]any{
			"name":        method.Name,
			"description": method.Description,
			"recipient":   method.Recipient,
		}).Error
	if err != nil {
		return errWrap(err)
	}

	return nil
}

func (r *RepositoryImpl) Delete(ctx context.Context, id int64) error {
	errWrap := getErrWrapper("Delete")

	err := r.db.WithContext(ctx).Delete(&paymentMethodEntity{}, "id = ?", id).Error
	if err != nil {
		return errWrap(err)
	}

	return nil
}

func getErrWrapper(method string) func(error) error {
	return func(err error) error {
		return fmt.Errorf("PaymentMethodRepositoryImpl->%s: %w", method, err)
	}
}
