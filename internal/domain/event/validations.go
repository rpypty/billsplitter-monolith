package event

import (
	apperrors "billsplitter-monolith/internal/errors"
	"billsplitter-monolith/internal/transport/http/middleware"
	"context"
)

func ValidateEventAccessBySession(ctx context.Context, ev *Event) error {
	sessionUserID, err := middleware.ExtractUserFromContext(ctx)
	if err != nil {
		return nil
	}

	for _, m := range ev.Members {
		if m.UserID != nil && *m.UserID == sessionUserID {
			return nil
		}
	}

	return apperrors.ErrForbiden
}
