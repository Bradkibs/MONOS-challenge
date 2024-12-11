package models

import (
	"github.com/google/uuid"
	"time"
)

type Invoice struct {
	ID        uuid.UUID  `json:"id"`
	PaymentID uuid.UUID  `json:"payment_id"`
	IssueDate time.Time  `json:"issue_date"`
	DueDate   time.Time  `json:"due_date"`
	Status    string     `json:"status"`
	DeletedAt *time.Time `json:"deleted_at"`
}
