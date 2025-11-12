package bill

import domainbill "billsplitter-monolith/internal/domain/bill"

func fromDomainBill(b domainbill.Bill) Bill {
	participants := make([]Participant, 0, len(b.Participants))

	for _, p := range b.Participants {
		participants = append(participants, Participant{
			MemberID: p.MemberID,
			Amount:   p.Amount,
		})
	}

	return Bill{
		ID:              b.ID,
		EventID:         b.EventID,
		Name:            b.Name,
		CreatedByUserID: int64(b.CreatedBy),
		TotalAmount:     int64(b.TotalAmount),
		Currency:        string(b.Currency),
		SplitType:       string(b.SplitTypeID),
		Participants:    participants,
		CreatedAt:       b.CreatedAt,
		UpdatedAt:       b.UpdatedAt,
		DeletedAt:       b.DeletedAt,
	}
}
