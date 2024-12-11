package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"

	"github.com/Bradkibs/MONOS-challenge/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AddPayment(payment *models.Payment, pool *pgxpool.Pool) error {
	var subscriptionStatus, tier, businessID string
	var branchCount int

	err := pool.QueryRow(context.Background(), `
		SELECT status, tier, businessId 
		FROM subscriptions WHERE id = $1`, payment.SubscriptionID).
		Scan(&subscriptionStatus, &tier, &businessID)
	if err != nil {
		return errors.New("subscription does not exist")
	}
	if subscriptionStatus != "active" {
		return errors.New("cannot add payment to an inactive subscription")
	}

	err = pool.QueryRow(context.Background(), `
		SELECT COUNT(*) FROM branches WHERE businessId = $1`, businessID).
		Scan(&branchCount)
	if err != nil {
		return errors.New("could not fetch branch count")
	}

	basePrices := map[string]float64{"Starter": 1.0, "Pro": 3.0, "Enterprise": 5.0}
	basePrice, exists := basePrices[tier]
	if !exists {
		return errors.New("invalid subscription tier")
	}

	expectedAmount := basePrice
	if branchCount > 1 {
		expectedAmount += float64(branchCount)
	}

	payment.Status = "completed"
	if payment.Amount != expectedAmount {
		payment.Status = "partial"
	}

	_, err = pool.Exec(context.Background(), `
		INSERT INTO payments (id, subscriptionId, amount, date, status) 
		VALUES ($1, $2, $3, $4, $5)`,
		payment.ID, payment.SubscriptionID, payment.Amount, payment.Date, payment.Status)
	if err != nil {
		return errors.New("failed to add payment to the database")
	}

	return nil
}

func GetPayment(paymentID string, pool *pgxpool.Pool) (*models.Payment, error) {
	var payment models.Payment
	err := pool.QueryRow(context.Background(), `
		SELECT id, subscriptionId, amount, date, status FROM payments WHERE id = $1`, paymentID).
		Scan(&payment.ID, &payment.SubscriptionID, &payment.Amount, &payment.Date, &payment.Status)
	if err != nil {
		return nil, errors.New("payment not found")
	}
	return &payment, nil
}

func GetPaymentsBySubscription(subscriptionID string, pool *pgxpool.Pool) ([]models.Payment, error) {
	rows, err := pool.Query(context.Background(), `
		SELECT id, subscriptionId, amount, date, status FROM payments WHERE subscriptionId = $1`, subscriptionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		if err := rows.Scan(&payment.ID, &payment.SubscriptionID, &payment.Amount, &payment.Date, &payment.Status); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func UpdatePayment(payment *models.Payment, pool *pgxpool.Pool) error {
	cmdTag, err := pool.Exec(context.Background(), `
		UPDATE payments SET amount = $2, date = $3, status = $4 WHERE id = $1`,
		payment.ID, payment.Amount, payment.Date, payment.Status)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("payment not found")
	}
	return nil
}

func DeletePayment(paymentID uuid.UUID, pool *pgxpool.Pool) error {
	query := `UPDATE payments SET deleted_at = NOW() WHERE id = $1`
	cmdTag, err := pool.Exec(context.Background(), query, paymentID)
	if err != nil {
		return fmt.Errorf("failed to delete payment: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("no rows were deleted, payment not found")
	}
	return nil
}

func HandleOverduePayment(subscriptionID string, pool *pgxpool.Pool) error {
	var endDate time.Time
	var status string
	err := pool.QueryRow(context.Background(), `
		SELECT endDate, status FROM subscriptions WHERE id = $1`, subscriptionID).
		Scan(&endDate, &status)
	if err != nil || status != "active" {
		return errors.New("invalid or inactive subscription")
	}

	if time.Now().After(endDate.Add(7 * 24 * time.Hour)) {
		_, err := pool.Exec(context.Background(), `
			UPDATE subscriptions SET status = 'suspended' WHERE id = $1`, subscriptionID)
		if err != nil {
			return err
		}
	}
	return nil
}

func HandlePartialPayment(paymentID string, pool *pgxpool.Pool) error {
	var amount float64
	var status string
	err := pool.QueryRow(context.Background(), `
		SELECT amount, status FROM payments WHERE id = $1`, paymentID).
		Scan(&amount, &status)
	if err != nil || status != "partial" {
		return errors.New("invalid payment or not partial")
	}

	_, err = pool.Exec(context.Background(), `
		UPDATE payments SET status = 'rejected' WHERE id = $1`, paymentID)
	if err != nil {
		return err
	}
	return errors.New("partial payment rejected, please retry with sufficient funds")
}
