package errors

import (
	"errors"
	"fmt"
)

var (
	ErrSessionNotFound        = errors.New("session not found")
	ErrUserNotFound           = errors.New("user not found")
	ErrEventNotFound          = errors.New("event not found")
	ErrFailedToGetUserFromCtx = errors.New("failed to get user from context")
	ErrSessionExpired         = errors.New("session expired")

	ErrValidationFunc = func(msg string) error {
		return fmt.Errorf("validation error: %s", msg)
	}
)
