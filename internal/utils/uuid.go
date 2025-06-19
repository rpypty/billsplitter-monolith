package utils

import (
	"github.com/google/uuid"
)

func NewUUIDv7() string {
	v, _ := uuid.NewV7()
	return v.String()
}
