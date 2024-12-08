package services

import (
	"context"
	"errors"
	"github.com/Bradkibs/MONOS-challenge/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func AddPayment(payment *models.Payment, pool *pgxpool.Pool) error {
	// Validate if the subscription exists and is active
	var subscriptionStatus string
	query := `SELECT status FROM subscriptions WHERE id = $1`
	err := pool.QueryRow(context.Background(), query, payment.SubscriptionID).Scan(&subscriptionStatus)
	if err != nil {
		return errors.New("subscription does not exist")
	}
	if subscriptionStatus != "active" {
		return errors.New("cannot add payment to an inactive subscription")
	}

	// Insert the payment into the database
	insertQuery := `INSERT INTO payments (id, subscriptionId, amount, date, status) VALUES ($1, $2, $3, $4, $5)`
	_, err = pool.Exec(context.Background(), insertQuery, payment.ID, payment.SubscriptionID, payment.Amount, payment.Date, payment.Status)
	if err != nil {
		return err
	}

	return nil
}

func GetPayment(paymentID string, pool *pgxpool.Pool) (*models.Payment, error) {
	query := `SELECT id, subscriptionId, amount, date, status FROM payments WHERE id = $1`
	var payment models.Payment
	err := pool.QueryRow(context.Background(), query, paymentID).Scan(&payment.ID, &payment.SubscriptionID, &payment.Amount, &payment.Date, &payment.Status)
	if err != nil {
		return nil, errors.New("payment not found")
	}
	return &payment, nil
}

func GetPaymentsBySubscription(subscriptionID string, pool *pgxpool.Pool) ([]models.Payment, error) {
	query := `SELECT id, subscriptionId, amount, date, status FROM payments WHERE subscriptionId = $1`
	rows, err := pool.Query(context.Background(), query, subscriptionID)
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
	query := `UPDATE payments SET amount = $2, date = $3, status = $4 WHERE id = $1`
	cmdTag, err := pool.Exec(context.Background(), query, payment.ID, payment.Amount, payment.Date, payment.Status)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("no rows were updated, payment not found")
	}

	return nil
}

func DeletePayment(paymentID string, pool *pgxpool.Pool) error {
	query := `DELETE FROM payments WHERE id = $1`
	cmdTag, err := pool.Exec(context.Background(), query, paymentID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("no rows were deleted, payment not found")
	}

	return nil
}

func HandleOverduePayment(subscriptionID string, pool *pgxpool.Pool) error {
	query := `SELECT endDate, status FROM subscriptions WHERE id = $1`
	var endDate time.Time
	var status string

	// Fetch subscription details
	err := pool.QueryRow(context.Background(), query, subscriptionID).Scan(&endDate, &status)
	if err != nil {
		return err
	}

	// Ensure the subscription is active
	if status != "active" {
		return errors.New("subscription is not active, cannot process overdue payment")
	}

	// Check if the subscription is overdue
	gracePeriod := 7 * 24 * time.Hour
	if time.Now().After(endDate.Add(gracePeriod)) {
		// Suspend the subscription
		updateQuery := `UPDATE subscriptions SET status = 'suspended' WHERE id = $1`
		_, err = pool.Exec(context.Background(), updateQuery, subscriptionID)
		if err != nil {
			return err
		}
	}

	return nil
}

func HandlePartialPayment(paymentID string, pool *pgxpool.Pool) error {
	query := `SELECT amount, status FROM payments WHERE id = $1`
	var amount float64
	var status string

	// Fetch payment details
	err := pool.QueryRow(context.Background(), query, paymentID).Scan(&amount, &status)
	if err != nil {
		return err
	}

	// Check if the payment is partial
	if status == "partial" {
		// Reject the payment
		updateQuery := `UPDATE payments SET status = 'rejected' WHERE id = $1`
		_, err = pool.Exec(context.Background(), updateQuery, paymentID)
		if err != nil {
			return err
		}

		return errors.New("partial payment rejected, please retry with sufficient funds")
	}

	return nil
}
