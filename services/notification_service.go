package services

import (
	"context"
	"fmt"
	"github.com/Bradkibs/MONOS-challenge/models"
	"github.com/Bradkibs/MONOS-challenge/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

func CreateNotification(pool *pgxpool.Pool, notification *models.Notification) error {
	query := `
		INSERT INTO notifications (userId, invoiceId, type, message, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, createdAt, updatedAt
	`
	err := pool.QueryRow(
		context.Background(),
		query,
		notification.UserID,
		notification.InvoiceID,
		notification.Type,
		notification.Message,
		notification.Status,
	).Scan(&notification.ID, &notification.CreatedAt, &notification.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	return nil
}
func GetNotificationByID(pool *pgxpool.Pool, id int) (*models.Notification, error) {
	query := `
		SELECT id, userId, invoiceId, type, message, status, createdAt, updatedAt
		FROM notifications
		WHERE id = $1
	`
	notification := &models.Notification{}
	err := pool.QueryRow(context.Background(), query, id).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.InvoiceID,
		&notification.Type,
		&notification.Message,
		&notification.Status,
		&notification.CreatedAt,
		&notification.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve notification by ID: %w", err)
	}

	return notification, nil
}
func UpdateNotification(pool *pgxpool.Pool, notification *models.Notification) error {
	query := `
		UPDATE notifications
		SET userId = $1, invoiceId = $2, type = $3, message = $4, status = $5, updatedAt = NOW()
		WHERE id = $6
	`
	_, err := pool.Exec(
		context.Background(),
		query,
		notification.UserID,
		notification.InvoiceID,
		notification.Type,
		notification.Message,
		notification.Status,
		notification.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}

	return nil
}
func DeleteNotification(pool *pgxpool.Pool, id int) error {
	query := `DELETE FROM notifications WHERE id = $1`
	_, err := pool.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	return nil
}
func SendReminderNotification(pool *pgxpool.Pool) error {
	// Query to fetch invoices with due dates within the next 3 days and status is not "paid"
	query := `
		SELECT i.id, i.dueDate, p.amount, u.id AS userId, u.email
		FROM invoices i
		JOIN payments p ON i.paymentId = p.id
		JOIN subscriptions s ON p.subscriptionId = s.id
		JOIN businesses b ON s.businessId = b.id
		JOIN users u ON b.vendorId = u.id
		WHERE i.status != 'paid' AND i.dueDate <= NOW() + INTERVAL '3 days'
	`
	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to query invoices for reminders: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var invoiceID, userID, email string
		var dueDate time.Time
		var amount float64

		err := rows.Scan(&invoiceID, &dueDate, &amount, &userID, &email)
		if err != nil {
			return fmt.Errorf("failed to scan invoice data: %w", err)
		}

		// Compose the notification message
		message := fmt.Sprintf(
			"Reminder: Your payment of $%.2f is due on %s. Please complete the payment to avoid penalties.",
			amount, dueDate.Format("2006-01-02"),
		)

		// Send the email notification
		err = utils.SendEmail(email, "Payment Reminder", message)
		if err != nil {
			log.Printf("failed to send reminder to %s: %v", email, err)
		}

		// Log the notification in the database
		err = CreateNotification(pool, &models.Notification{
			UserID:    userID,
			InvoiceID: invoiceID,
			Type:      "Reminder",
			Message:   message,
			Status:    "pending", // Assume "pending" until email service confirms sending
		})
		if err != nil {
			log.Printf("failed to log notification for user %s: %v", userID, err)
		}
	}

	if rows.Err() != nil {
		return fmt.Errorf("error iterating over invoice rows: %w", rows.Err())
	}

	return nil
}
func SendLatePaymentAlerts(pool *pgxpool.Pool) error {
	// Query to fetch invoices with overdue payments (status is not "paid" and due date has passed)
	query := `
		SELECT i.id, i.dueDate, p.amount, u.id AS userId, u.email
		FROM invoices i
		JOIN payments p ON i.paymentId = p.id
		JOIN subscriptions s ON p.subscriptionId = s.id
		JOIN businesses b ON s.businessId = b.id
		JOIN users u ON b.vendorId = u.id
		WHERE i.status != 'paid' AND i.dueDate < NOW()
	`
	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to query overdue invoices: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var invoiceID, userID, email string
		var dueDate time.Time
		var amount float64

		err := rows.Scan(&invoiceID, &dueDate, &amount, &userID, &email)
		if err != nil {
			return fmt.Errorf("failed to scan overdue invoice data: %w", err)
		}

		// Compose the notification message
		message := fmt.Sprintf(
			"Alert: Your payment of $%.2f was due on %s and is now overdue. Please make the payment immediately to avoid further penalties.",
			amount, dueDate.Format("2006-01-02"),
		)

		// Send the email notification
		err = utils.SendEmail(email, "Overdue Payment Alert", message)
		if err != nil {
			log.Printf("failed to send overdue alert to %s: %v", email, err)
		}

		// Log the notification in the database
		err = CreateNotification(pool, &models.Notification{
			UserID:    userID,
			InvoiceID: invoiceID,
			Type:      "LatePayment",
			Message:   message,
			Status:    "pending", // Assume "pending" until email service confirms sending
		})
		if err != nil {
			log.Printf("failed to log notification for user %s: %v", userID, err)
		}
	}

	if rows.Err() != nil {
		return fmt.Errorf("error iterating over overdue invoice rows: %w", rows.Err())
	}

	return nil
}
