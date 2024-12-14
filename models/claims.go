package models

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type Claims struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	Email     string     `json:"email"`
	IssuedAt  time.Time  `json:"issued_at"`
	ExpiresAt time.Time  `json:"expires_at"`
	Role      string     `json:"role"`
	Status    string     `json:"status"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (c Claims) Valid() error {
	if c.ExpiresAt.Before(time.Now()) {
		return errors.New("token has expired")
	}

	if c.DeletedAt != nil {
		return errors.New("account is deactivated or deleted")
	}

	return nil
}
