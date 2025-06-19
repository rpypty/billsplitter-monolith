package errors

import (
	"errors"
	"fmt"
)

var (
	ErrSessionNotFound        = errors.New("session not found")
	ErrFailedToGetUserFromCtx = errors.New("failed to get user from context")

	ErrUserStorageFunc = func(err error, method string) error {
		return fmt.Errorf("USER_STORAGE_ERROR->%s(): %w", method, err)
	}
	ErrUserServiceFunc = func(err error, method string) error {
		return fmt.Errorf("USER_SERVICE_ERROR->%s(): %w", method, err)
	}
	ErrSessionStorageFunc = func(err error, method string) error {
		return fmt.Errorf("SESSION_STORAGE_ERROR->%s(): %w", method, err)
	}
)
