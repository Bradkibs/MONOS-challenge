package models

import "github.com/google/uuid"

type Branch struct {
	ID         uuid.UUID `json:"id"`
	BusinessID string    `json:"business_id"`
	Location   string    `json:"location"`
}
