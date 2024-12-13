package models

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ID         uuid.UUID  `json:"id"`
	BusinessID uuid.UUID  `json:"business_id"`
	Name       string     `json:"name"`
	Details    string     `json:"details"`
	Quantity   int        `json:"quantity"`
	Price      float64    `json:"price"`
	DeletedAt  *time.Time `json:"deleted_at"`
}
