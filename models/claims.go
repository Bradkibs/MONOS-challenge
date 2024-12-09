package models

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
