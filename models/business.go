package models

import (
	"github.com/google/uuid"
	"time"
)

type Business struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	VendorID    uuid.UUID  `json:"vendor_id"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
