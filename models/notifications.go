package models

import "time"

type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	InvoiceID string    `json:"invoiceId"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
