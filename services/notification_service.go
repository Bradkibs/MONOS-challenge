package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Bradkibs/MONOS-challenge/models"
	"github.com/Bradkibs/MONOS-challenge/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateNotification(pool *pgxpool.Pool, notification *models.Notification) error {
	query := `
		INSERT INTO notifications (user_id, invoice_id, type, message, createdat, updatedat, deletedat)
		VALUES ($1, $2, $3, $4, NOW(), NOW(), NULL)
		RETURNING id, createdat, updatedat
	`
	err := pool.QueryRow(
		context.Background(),
		query,
		notification.UserID,
		notification.InvoiceID,
		notification.Type,
		notification.Message,
	).Scan(&notification.ID, &notification.CreatedAt, &notification.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	return nil
}

func GetNotificationByID(pool *pgxpool.Pool, id uuid.UUID) (*models.Notification, error) {
	query := `
		SELECT id, user_id, invoice_id, type, message, createdat, updatedat, deletedat
		FROM notifications WHERE id = $1 AND deletedat IS NULL
	`
	notification := &models.Notification{}
	err := pool.QueryRow(context.Background(), query, id).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.InvoiceID,
		&notification.Type,
		&notification.Message,
		&notification.CreatedAt,
		&notification.UpdatedAt,
		&notification.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve notification by ID: %w", err)
	}
	return notification, nil
}

func GetNotificationByUserID(pool *pgxpool.Pool, userId uuid.UUID) (*models.Notification, error) {
	query := `
		SELECT id, user_id, invoice_id, type, message, createdat, updatedat, deletedat
		FROM notifications WHERE user_id = $1 AND deletedat IS NULL
	`
	notification := &models.Notification{}
	err := pool.QueryRow(context.Background(), query, userId).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.InvoiceID,
		&notification.Type,
		&notification.Message,
		&notification.CreatedAt,
		&notification.UpdatedAt,
		&notification.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve notification by ID: %w", err)
	}
	return notification, nil
}
func GetNotificationByInvoiceID(pool *pgxpool.Pool, invoiceId uuid.UUID) (*models.Notification, error) {
	query := `
		SELECT id, user_id, invoice_id, type, message, createdat, updatedat, deletedat
		FROM notifications WHERE invoice_id = $1 AND deletedat IS NULL
	`
	notification := &models.Notification{}
	err := pool.QueryRow(context.Background(), query, invoiceId).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.InvoiceID,
		&notification.Type,
		&notification.Message,
		&notification.CreatedAt,
		&notification.UpdatedAt,
		&notification.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve notification by ID: %w", err)
	}
	return notification, nil
}

func UpdateNotification(pool *pgxpool.Pool, notification *models.Notification) error {
	query := `
		UPDATE notifications SET type = $1, message = $2, updatedat = NOW() WHERE id = $3 AND deletedat IS NULL
	`
	cmdTag, err := pool.Exec(
		context.Background(),
		query,
		notification.Type,
		notification.Message,
		notification.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("no rows were updated, notification not found")
	}
	return nil
}

func DeleteNotification(pool *pgxpool.Pool, id uuid.UUID) error {
	query := `UPDATE notifications SET deletedat = NOW() WHERE id = $1 AND deletedat IS NULL`
	cmdTag, err := pool.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("no rows were deleted, notification not found")
	}
	return nil
}

func SendReminderNotification(pool *pgxpool.Pool) error {
	query := `
		SELECT i.id, i.duedate, p.amount, u.id AS user_id, u.email
		FROM invoices i
		JOIN payments p ON i.paymentid = p.id
		JOIN subscriptions s ON p.subscriptionid = s.id
		JOIN businesses b ON s.businessid = b.id
		JOIN users u ON b.vendorid = u.id
		WHERE i.status != 'paid' AND i.duedate <= NOW() + INTERVAL '3 days'
		AND i.deletedat IS NULL
	`
	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to query invoices for reminders: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var invoiceID, userID uuid.UUID
		var email string
		var dueDate time.Time
		var amount float64

		if err := rows.Scan(&invoiceID, &dueDate, &amount, &userID, &email); err != nil {
			log.Printf("failed to scan invoice data: %v", err)
			continue
		}

		message := fmt.Sprintf("Reminder: Your payment of $%.2f is due on %s.", amount, dueDate.Format("2006-01-02"))
		if err := utils.SendEmail(email, "Payment Reminder", message); err != nil {
			log.Printf("failed to send reminder to %s: %v", email, err)
		}

		if err := CreateNotification(pool, &models.Notification{
			UserID:    userID,
			InvoiceID: &invoiceID,
			Type:      "Reminder",
			Message:   message,
		}); err != nil {
			log.Printf("failed to log notification for user %s: %v", userID, err)
		}
	}

	return rows.Err()
}
