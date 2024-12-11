package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/Bradkibs/MONOS-challenge/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"os"
	"regexp"
	"time"
	"unicode"
)

var jwtSecret = os.Getenv("JWT_SECRET")

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

func GenerateJWT(userID, emailOrPhoneNumber, role string, deletedAt *time.Time) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		ID:        uuid.New(),
		UserID:    uuid.MustParse(userID),
		Email:     emailOrPhoneNumber,
		IssuedAt:  time.Now(),
		ExpiresAt: expirationTime,
		Role:      role,
		DeletedAt: deletedAt,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ParseJWT(tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}
	return claims, nil
}

func RegisterUserByEmail(email, password, role string, pool *pgxpool.Pool) (string, error) {
	if !isValidEmail(email) {
		return "", errors.New("invalid email format")
	}
	if !isValidPassword(password) {
		return "", errors.New("password must be at least 8 characters long and contain a mix of letters, numbers, and special characters")
	}

	var existingUserID uuid.UUID
	err := pool.QueryRow(context.Background(), "SELECT id FROM users WHERE email = $1", email).Scan(&existingUserID)
	if err == nil {
		return "", errors.New("user with this email already exists")
	} else if err != pgx.ErrNoRows {
		return "", fmt.Errorf("failed to check for existing user: %v", err)
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}

	var userID uuid.UUID
	err = pool.QueryRow(context.Background(), "INSERT INTO users (email, password, role) VALUES ($1, $2, $3) RETURNING id", email, hashedPassword, role).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %v", err)
	}

	token, err := GenerateJWT(userID.String(), email, role, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}
	return token, nil
}
func RegisterUserByPhoneNumber(phoneNumber, password, role string, pool *pgxpool.Pool) (string, error) {
	if !isValidPhoneNumber(phoneNumber) {
		return "", errors.New("invalid Phone number format")
	}
	if !isValidPassword(password) {
		return "", errors.New("password must be at least 8 characters long and contain a mix of letters, numbers, and special characters")
	}

	var existingUserID uuid.UUID
	err := pool.QueryRow(context.Background(), "SELECT id FROM users WHERE email = $1", phoneNumber).Scan(&existingUserID)
	if err == nil {
		return "", errors.New("user with this email already exists")
	} else if err != pgx.ErrNoRows {
		return "", fmt.Errorf("failed to check for existing user: %v", err)
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}

	var userID uuid.UUID
	err = pool.QueryRow(context.Background(), "INSERT INTO users (phone_number, password, role) VALUES ($1, $2, $3) RETURNING id", phoneNumber, hashedPassword, role).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %v", err)
	}

	token, err := GenerateJWT(userID.String(), phoneNumber, role, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}
	return token, nil
}

func isValidPhoneNumber(phone string) bool {
	re := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	return re.MatchString(phone)
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasNumber, hasSpecial, hasLetter := false, false, false
	for _, char := range password {
		switch {
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		case unicode.IsLetter(char):
			hasLetter = true
		}
	}
	return hasNumber && hasSpecial && hasLetter
}

func LoginUser(email, password string, pool *pgxpool.Pool) (string, error) {
	var userID uuid.UUID
	var hashedPassword, role string
	var deletedAt *time.Time
	err := pool.QueryRow(context.Background(), "SELECT id, password, role, deleted_at FROM users WHERE email = $1", email).Scan(&userID, &hashedPassword, &role, &deletedAt)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if deletedAt != nil {
		return "", errors.New("account is deactivated")
	}

	if err := CheckPassword(hashedPassword, password); err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := GenerateJWT(userID.String(), email, role, deletedAt)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}
	return token, nil
}
