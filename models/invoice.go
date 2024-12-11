package models

import (
	"github.com/google/uuid"
	"time"
)

type Invoice struct {
	ID        uuid.UUID  `json:"id"`         // UUID as primary key
	PaymentID uuid.UUID  `json:"payment_id"` // Foreign key referencing Payments
	IssueDate time.Time  `json:"issue_date"` // Issue date as DATE type
	DueDate   time.Time  `json:"due_date"`   // Due date as DATE type
	Status    string     `json:"status"`     // Status of the invoice
	DeletedAt *time.Time `json:"deleted_at"` // Soft delete column, nullable
}
