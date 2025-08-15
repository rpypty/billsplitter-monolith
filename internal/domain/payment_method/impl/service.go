package impl

import (
	"context"

	domain "billsplitter-monolith/internal/domain/payment_method"
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

func (s *ServiceImpl) FetchByUserID(ctx context.Context, userID int64) ([]domain.PaymentMethod, error) {
	return s.repo.FetchByUserID(ctx, userID)
}

func (s *ServiceImpl) Create(ctx context.Context, paymentMethod domain.PaymentMethod) (domain.PaymentMethod, error) {
	return s.repo.Create(ctx, paymentMethod)
}

func (s *ServiceImpl) Update(ctx context.Context, id int64, paymentMethod domain.PaymentMethod) error {
	return s.repo.Update(ctx, id, paymentMethod)
}

func (s *ServiceImpl) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
