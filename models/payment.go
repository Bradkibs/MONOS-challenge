package models

import (
	"github.com/google/uuid"
	"time"
)

type Payment struct {
	ID             uuid.UUID  `json:"id"`
	SubscriptionID uuid.UUID  `json:"subscription_id"`
	Amount         float64    `json:"amount"`
	Date           time.Time  `json:"date"`
	Status         string     `json:"status"`
	DeletedAt      *time.Time `json:"deleted_at"`
}
