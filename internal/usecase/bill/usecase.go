package bill

import (
	"context"
	"fmt"

	domainbill "billsplitter-monolith/internal/domain/bill"
	domainevent "billsplitter-monolith/internal/domain/event"
	vo "billsplitter-monolith/internal/domain/valueobject"
	"billsplitter-monolith/internal/errors"
)

type UseCase interface {
	CreateBill(ctx context.Context, rq CreateBillRq) (int64, error)
	FetchEventBills(ctx context.Context, eventID int64) ([]domainbill.Bill, error)
}

type UseCaseImpl struct {
	billSvc  domainbill.Service
	eventSvc domainevent.Service
}

type CreateBillRq struct {
	EventID      int64
	Name         string
	CreatedBy    vo.UserID // userID из сессии инициатора запроса на создание чека
	TotalAmount  vo.Amount
	Currency     vo.CurrencyCode
	SplitType    vo.SplitType
	PaidBy       int64 // memberID - не путать с userID. memberID - это айди сущности участника ивента, а не профиля пользователя - это разные таблицы (user и event_members)
	Participants []domainbill.Participant
}

func New(billSvc domainbill.Service, eventSvc domainevent.Service) *UseCaseImpl {
	return &UseCaseImpl{
		billSvc:  billSvc,
		eventSvc: eventSvc,
	}
}

func (uc *UseCaseImpl) CreateBill(ctx context.Context, rq CreateBillRq) (int64, error) {
	eventModel, err := uc.getEvent(ctx, rq.EventID)
	if err != nil {
		return 0, err
	}

	if err := uc.validateRequest(rq); err != nil {
		return 0, err
	}

	if err := uc.validateParticipants(eventModel, rq); err != nil {
		return 0, err
	}

	billID, err := uc.billSvc.Create(ctx, domainbill.Bill{
		Name:         rq.Name,
		CreatedBy:    rq.CreatedBy,
		Participants: rq.Participants,
		PaidBy:       rq.PaidBy,
		EventID:      rq.EventID,
		TotalAmount:  rq.TotalAmount,
		Currency:     rq.Currency,
		SplitTypeID:  rq.SplitType,
	})
	if err != nil {
		return 0, err
	}

	return billID, nil
}

func (uc *UseCaseImpl) FetchEventBills(ctx context.Context, eventID int64) ([]domainbill.Bill, error) {
	// достаем ради валидации
	_, err := uc.getEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return uc.billSvc.FetchByEventID(ctx, eventID)
}

func (uc *UseCaseImpl) getEvent(ctx context.Context, eventID int64) (*domainevent.Event, error) {
	ev, err := uc.eventSvc.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if ev == nil || ev.ID == 0 {
		return nil, errors.ErrEventNotFound
	}

	if err := domainevent.ValidateEventAccessBySession(ctx, ev); err != nil {
		return nil, err
	}

	return ev, nil
}

func (uc *UseCaseImpl) validateParticipants(ev *domainevent.Event, rq CreateBillRq) error {
	if rq.PaidBy == 0 {
		return errors.ErrValidationFunc("paid_by must be provided")
	}

	if rq.CreatedBy == 0 {
		return errors.ErrValidationFunc("created_by must be provided")
	}

	var createdByBelogsToEvent bool

	memberIndex := make(map[int64]struct{}, len(ev.Members))
	for _, m := range ev.Members {
		memberIndex[m.ID] = struct{}{}

		if m.UserID != nil && *m.UserID == rq.CreatedBy {
			createdByBelogsToEvent = true
		}
	}

	// проверяем что инициатор запроса вообще принадлежит этому инвета
	if !createdByBelogsToEvent {
		return errors.ErrValidationFunc(fmt.Sprintf("created_by %d does not belong to event %d", rq.CreatedBy, ev.ID))
	}

	// проверяем что юзер который платит принадлежит этому ивенту
	if _, ok := memberIndex[int64(rq.PaidBy)]; !ok {
		return errors.ErrValidationFunc(fmt.Sprintf("paid_by %d does not belong to event %d", rq.PaidBy, ev.ID))
	}

	// проверяем что все юзеры в чеке принадлежат этому ивенту
	for _, p := range rq.Participants {
		if _, ok := memberIndex[p.MemberID]; !ok {
			return errors.ErrValidationFunc(fmt.Sprintf("member_id %d does not belong to event %d", p.MemberID, ev.ID))
		}
	}

	return nil
}

func (uc *UseCaseImpl) validateRequest(rq CreateBillRq) error {
	if rq.TotalAmount <= 0 {
		return errors.ErrValidationFunc(fmt.Sprintf("bill amount must be positive number"))
	}

	var amountSum int64
	for _, participant := range rq.Participants {
		amountSum += participant.Amount
	}

	if amountSum != int64(rq.TotalAmount) {
		return errors.ErrValidationFunc(fmt.Sprintf("bill amount (%d) not equal to sum of shares (%d)", rq.TotalAmount, amountSum))
	}

	return nil
}
