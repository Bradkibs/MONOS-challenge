package utils

import (
	"github.com/google/uuid"
)

// GenerateUniqueID generates a new unique UUID string.
func GenerateUniqueID() string {
	return uuid.New().String()
}
