package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"regexp"
	"unicode"

	"github.com/Bradkibs/MONOS-challenge/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"os"

	"time"
)

// JWT secret key (use environment variables for better security)
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

// GenerateJWT generates a new JWT token for the user
func GenerateJWT(userID, email string) (string, error) {
	// Set expiration time
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the claims
	claims := &models.Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ParseJWT(tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	// Check if the token is valid
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	return claims, nil
}

func RegisterUser(email, password string, pool *pgxpool.Pool) (string, error) {
	// Validate the email format
	if !isValidEmail(email) {
		return "", errors.New("invalid email format")
	}

	// Validate the password strength
	if !isValidPassword(password) {
		return "", errors.New("password must be at least 8 characters long and contain a mix of letters, numbers, and special characters")
	}

	// Check if the email already exists
	existingUserQuery := `SELECT id FROM users WHERE email = $1`
	var existingUserID string
	err := pool.QueryRow(context.Background(), existingUserQuery, email).Scan(&existingUserID)
	if err == nil {
		return "", errors.New("user with this email already exists")
	} else if err != pgx.ErrNoRows {
		// Handle unexpected database errors
		return "", fmt.Errorf("failed to check for existing user: %v", err)
	}

	// Hash the password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}

	// Save the user to the database
	createUserQuery := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`
	var userID string
	err = pool.QueryRow(context.Background(), createUserQuery, email, hashedPassword).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %v", err)
	}

	// Generate a JWT for the new user
	token, err := GenerateJWT(userID, email)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return token, nil
}

func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasNumber := false
	hasSpecial := false
	hasLetter := false
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
	// Retrieve user from the database
	query := `SELECT id, password FROM users WHERE email = $1`
	var userID, hashedPassword string
	err := pool.QueryRow(context.Background(), query, email).Scan(&userID, &hashedPassword)
	if err != nil {
		return "", fmt.Errorf("invalid email or password")
	}

	// Check the password
	err = CheckPassword(hashedPassword, password)
	if err != nil {
		return "", fmt.Errorf("invalid email or password")
	}

	// Generate a JWT for the user
	token, err := GenerateJWT(userID, email)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return token, nil
}
