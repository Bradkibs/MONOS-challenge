package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	PhoneNumber string     `json:"phone_number"`
	Email       string     `json:"email"`
	Password    string     `json:"password"`
	Role        string     `json:"role"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
