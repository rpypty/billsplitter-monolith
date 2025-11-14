package meet

import (
	domainevent "billsplitter-monolith/internal/domain/event"
	"billsplitter-monolith/internal/usecase/event"
	"billsplitter-monolith/internal/utils"
)

func fromDomainEvent(meet domainevent.Event) Event {
	outMembers := make([]Member, 0, len(meet.Members))

	for _, member := range meet.Members {
		outMembers = append(outMembers, Member{
			ID:       member.ID,
			UserID:   utils.UserIDToInt64(member.UserID),
			Username: member.Name,
		})
	}

	return Event{
		ID:              meet.ID,
		Name:            meet.Name,
		CreatedByUserID: int64(meet.CreatedByUserID),
		Members:         outMembers,
		Status:          string(meet.Status),
		Type:            string(meet.Type),
		CreatedAt:       meet.CreatedAt,
		UpdatedAt:       meet.UpdatedAt,
		EventDate:       meet.EventDate,
		DeletedAt:       meet.DeletedAt,
	}
}

func fromDomainSummary(summary event.EventSummary) EventSummary {
	outBalances := make([]Balance, 0, len(summary.Balances))
	for _, balance := range summary.Balances {
		outBalances = append(outBalances, Balance{
			MemberID:   balance.MemberID,
			UserID:     utils.UserIDToInt64(balance.UserID),
			Name:       balance.Name,
			TotalPaid:  balance.TotalPaid,
			TotalShare: balance.TotalShare,
			Balance:    balance.Balance,
		})
	}

	outSettlements := make([]Settlement, 0, len(summary.Settlements))
	for _, settlement := range summary.Settlements {
		outSettlements = append(outSettlements, Settlement{
			FromMemberID: settlement.FromMemberID,
			ToMemberID:   settlement.ToMemberID,
			Amount:       settlement.Amount,
		})
	}

	return EventSummary{
		Balances:    outBalances,
		Settlements: outSettlements,
	}
}
