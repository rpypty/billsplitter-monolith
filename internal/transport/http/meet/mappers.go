package meet

import (
	domain "billsplitter-monolith/internal/domain/event"
	"billsplitter-monolith/internal/utils"
)

func fromDomainEvent(meet domain.Event) Event {
	outMembers := make([]Member, 0, len(meet.Members))

	for _, member := range meet.Members {
		outMembers = append(outMembers, Member{
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
