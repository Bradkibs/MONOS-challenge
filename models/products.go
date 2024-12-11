package models

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ID         uuid.UUID  `json:"id"`          // UUID as primary key
	BusinessID uuid.UUID  `json:"business_id"` // Foreign key referencing Businesses
	Name       string     `json:"name"`        // Product name, non-null
	Details    string     `json:"details"`     // Optional product details (TEXT)
	Quantity   int        `json:"quantity"`    // Quantity, non-null
	Price      float64    `json:"price"`       // Price as DECIMAL, non-null
	DeletedAt  *time.Time `json:"deleted_at"`  // Soft delete column, nullable
}
