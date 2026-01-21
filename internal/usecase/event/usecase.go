package event

import (
	"context"
	"fmt"
	"sort"
	"time"

	"billsplitter-monolith/internal/domain/bill"
	"billsplitter-monolith/internal/domain/event"
	"billsplitter-monolith/internal/domain/user"
	vo "billsplitter-monolith/internal/domain/valueobject"
	"billsplitter-monolith/internal/errors"
	"billsplitter-monolith/internal/utils"
)

type UseCase interface {
	CreateMeet(ctx context.Context, rq CreateMeetRq) (int64, error)
	FetchUserMeets(ctx context.Context, userID vo.UserID) ([]event.Event, error)
	GetMeetByID(ctx context.Context, meetID int64) (*event.Event, error)
	CalculateSummary(ctx context.Context, meetID int64) (*EventSummary, error)
	AssignMemberToUser(ctx context.Context, meetID int64, memberID int64, userID vo.UserID) error
}

type UseCaseImpl struct {
	eventSvc event.Service
	userSvc  user.Service
	billSvc  bill.Service
}

func New(
	eventSvc event.Service,
	userSvc user.Service,
	billSvc bill.Service,
) *UseCaseImpl {
	return &UseCaseImpl{
		eventSvc: eventSvc,
		userSvc:  userSvc,
		billSvc:  billSvc,
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

	if rq.Date == nil {
		rq.Date = utils.Ptr(time.Now())
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
	ev, err := uc.eventSvc.GetByID(ctx, meetID)
	if err != nil {
		return nil, err
	}

	if ev == nil || ev.ID == 0 {
		return nil, errors.ErrEventNotFound
	}

	if err := event.ValidateEventAccessBySession(ctx, ev); err != nil {
		return nil, err
	}

	return ev, nil
}

func (uc *UseCaseImpl) CalculateSummary(ctx context.Context, meetID int64) (*EventSummary, error) {
	meet, err := uc.GetMeetByID(ctx, meetID)
	if err != nil {
		return nil, err
	}

	if meet == nil {
		return nil, errors.ErrEventNotFound
	}

	bills, err := uc.billSvc.FetchByEventID(ctx, meet.ID)
	if err != nil {
		return nil, err
	}

	balances := buildBalances(*meet, bills)
	settlements := calculateSettlements(balances)

	return &EventSummary{
		Balances:    balances,
		Settlements: settlements,
	}, nil
}

func (uc *UseCaseImpl) AssignMemberToUser(ctx context.Context, meetID int64, memberID int64, userID vo.UserID) error {
	if meetID <= 0 {
		return errors.ErrValidationFunc("meet_id must be provided")
	}
	if memberID <= 0 {
		return errors.ErrValidationFunc("member_id must be provided")
	}

	meet, err := uc.eventSvc.GetByID(ctx, meetID)
	if err != nil {
		return err
	}
	if meet == nil || meet.ID == 0 {
		return errors.ErrEventNotFound
	}

	var selectedMember *event.Member
	for i := range meet.Members {
		if meet.Members[i].ID == memberID {
			selectedMember = &meet.Members[i]
			break
		}
	}
	if selectedMember == nil {
		return errors.ErrValidationFunc(fmt.Sprintf("member_id %d does not belong to event %d", memberID, meetID))
	}
	if selectedMember.UserID != nil {
		return errors.ErrValidationFunc("member already assigned")
	}

	return uc.eventSvc.AssignMemberUser(ctx, meetID, memberID, userID)
}

func buildBalances(ev event.Event, bills []bill.Bill) []Balance {
	totalPaid := make(map[int64]int64, len(ev.Members))
	totalShare := make(map[int64]int64, len(ev.Members))

	for _, b := range bills {
		totalPaid[b.PaidBy] += int64(b.TotalAmount)

		for _, participant := range b.Participants {
			totalShare[participant.MemberID] += participant.Amount
		}
	}

	balances := make([]Balance, 0, len(ev.Members))
	for _, member := range ev.Members {
		paid := totalPaid[member.ID]
		share := totalShare[member.ID]

		balances = append(balances, Balance{
			MemberID:   member.ID,
			UserID:     member.UserID,
			Name:       member.Name,
			TotalPaid:  paid,
			TotalShare: share,
			Balance:    paid - share,
		})
	}

	return balances
}

func calculateSettlements(balances []Balance) []Settlement {
	debtors := make([]settlementParticipant, 0)
	creditors := make([]settlementParticipant, 0)

	for _, balance := range balances {
		switch {
		case balance.Balance < 0:
			debtors = append(debtors, settlementParticipant{
				MemberID: balance.MemberID,
				Amount:   -balance.Balance,
			})
		case balance.Balance > 0:
			creditors = append(creditors, settlementParticipant{
				MemberID: balance.MemberID,
				Amount:   balance.Balance,
			})
		}
	}

	sort.Slice(debtors, func(i, j int) bool {
		if debtors[i].Amount == debtors[j].Amount {
			return debtors[i].MemberID < debtors[j].MemberID
		}
		return debtors[i].Amount < debtors[j].Amount
	})

	sort.Slice(creditors, func(i, j int) bool {
		if creditors[i].Amount == creditors[j].Amount {
			return creditors[i].MemberID < creditors[j].MemberID
		}
		return creditors[i].Amount < creditors[j].Amount
	})

	capacity := len(debtors)
	if len(creditors) < capacity {
		capacity = len(creditors)
	}

	settlements := make([]Settlement, 0, capacity)
	for i, j := 0, 0; i < len(debtors) && j < len(creditors); {
		amount := minInt64(debtors[i].Amount, creditors[j].Amount)
		settlements = append(settlements, Settlement{
			FromMemberID: debtors[i].MemberID,
			ToMemberID:   creditors[j].MemberID,
			Amount:       amount,
		})

		debtors[i].Amount -= amount
		creditors[j].Amount -= amount

		if debtors[i].Amount == 0 {
			i++
		}

		if creditors[j].Amount == 0 {
			j++
		}
	}

	return settlements
}

type settlementParticipant struct {
	MemberID int64
	Amount   int64
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
