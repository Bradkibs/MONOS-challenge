package models

import (
	"github.com/google/uuid"
	"time"
)

type Subscription struct {
	ID         uuid.UUID  `json:"id"`
	BusinessID uuid.UUID  `json:"business_id"`
	Tier       string     `json:"tier"`
	StartDate  time.Time  `json:"start_date"`
	EndDate    *time.Time `json:"end_date"`
	Status     string     `json:"status"`
	DeletedAt  *time.Time `json:"deleted_at"`
}
