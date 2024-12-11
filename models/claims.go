package models

import (
	"github.com/google/uuid"
	"time"
)

type Claims struct {
	ID        uuid.UUID  `json:"id"`         // UUID primary key
	UserID    uuid.UUID  `json:"user_id"`    // Foreign key referencing users
	Email     string     `json:"email"`      // Email associated with the claim
	IssuedAt  time.Time  `json:"issued_at"`  // Timestamp when the claim was issued
	ExpiresAt time.Time  `json:"expires_at"` // Timestamp when the claim expires
	Role      string     `json:"role"`       // Role of the user (nullable)
	Status    string     `json:"status"`     // Status of the claim, default 'active'
	DeletedAt *time.Time `json:"deleted_at"` // Soft delete column, nullable
}
