package impl

import (
	"context"

	"billsplitter-monolith/internal/domain/event"
	vo "billsplitter-monolith/internal/domain/valueobject"
	"billsplitter-monolith/internal/utils"
)

var _ event.Service = (*ServiceImpl)(nil)

type ServiceImpl struct {
	repo event.Repository
}

func New(repo event.Repository) *ServiceImpl {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) Create(ctx context.Context, rq event.CreateEventRq) (int64, error) {
	// TODO: add pg transaction

	billID, err := s.repo.Create(ctx, event.Event{
		Name:            rq.Name,
		PublicUUID:      utils.NewUUIDv4(),
		CreatedByUserID: rq.CreatorUserID,
		Status:          event.StatusDraft,
		Type:            event.TypeMeet,
		EventDate:       rq.EventDate,
	})
	if err != nil {
		return 0, err
	}

	// +1 чтобы добавить капасити для создателя ивента
	members := make([]event.Member, 0, len(rq.Members)+1)

	members = append(members, event.Member{
		UserID: &rq.CreatorUserID,
		Name:   rq.CreatorUsername,
	})

	for _, memberName := range rq.Members {
		members = append(members, event.Member{
			UserID: nil, // еще не знаем айди пользователя, он появится когда нажмет кнопку "пикми"
			Name:   memberName,
		})
	}

	err = s.repo.AddMembers(ctx, billID, members)
	if err != nil {
		return 0, err
	}

	return billID, nil
}

func (s *ServiceImpl) Update(ctx context.Context, rq event.UpdateEventRq) error {
	// TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) AddUsers(ctx context.Context, eventID int64, userID []vo.UserID) error {
	// TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) RemoveUsers(ctx context.Context, eventID int64, userID []vo.UserID) error {
	// TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) Delete(ctx context.Context, billID int64) error {
	// TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) FetchByUserID(ctx context.Context, userID vo.UserID) ([]event.Event, error) {
	return s.repo.FetchByUserID(ctx, userID)
}

func (s *ServiceImpl) GetByID(ctx context.Context, billID int64) (*event.Event, error) {
	return s.repo.GetByID(ctx, billID)
}

func (s *ServiceImpl) GetByPublicUUID(ctx context.Context, publicUUID string) (*event.Event, error) {
	return s.repo.GetByPublicUUID(ctx, publicUUID)
}

func (s *ServiceImpl) AssignMemberUser(ctx context.Context, eventID int64, memberID int64, userID vo.UserID) error {
	return s.repo.AssignMemberUser(ctx, eventID, memberID, userID)
}
