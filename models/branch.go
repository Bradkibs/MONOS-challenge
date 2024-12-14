package models

import "github.com/google/uuid"

type Branch struct {
	ID         uuid.UUID `json:"id"`
	BusinessID string    `json:"business_id"`
	Country    string    `json:"country"`
	Location   string    `json:"location"`
}
