package impl

import (
	"context"

	"billsplitter-monolith/internal/domain/bill"
	vo "billsplitter-monolith/internal/domain/valueobject"
)

var _ bill.Service = (*ServiceImpl)(nil)

type ServiceImpl struct {
	repo bill.Repository
}

func (s ServiceImpl) Create(ctx context.Context, bill bill.Bill) (int64, error) {
	return s.repo.Create(ctx, bill)
}

func (s ServiceImpl) FetchByEventID(ctx context.Context, eventID int64) ([]bill.Bill, error) {
	return s.repo.FetchByEventID(ctx, eventID)
}

func (s ServiceImpl) Update(ctx context.Context, rq bill.UpdateBillRq) error {
	// TODO implement me
	panic("implement me")
}

func (s ServiceImpl) AddUsers(ctx context.Context, billID int64, users []vo.UserID) error {
	// TODO implement me
	panic("implement me")
}

func (s ServiceImpl) RemoveUsers(ctx context.Context, billID int64, users []vo.UserID) error {
	// TODO implement me
	panic("implement me")
}

func (s ServiceImpl) Delete(ctx context.Context, billID int64) error {
	// TODO implement me
	panic("implement me")
}

func New(repo bill.Repository) *ServiceImpl {
	return &ServiceImpl{
		repo: repo,
	}
}
