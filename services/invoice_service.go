package services

import (
	"context"
	"errors"
	"github.com/Bradkibs/MONOS-challenge/models"
	"github.com/Bradkibs/MONOS-challenge/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func AddInvoice(invoice *models.Invoice, pool *pgxpool.Pool) error {
	// Check payment status
	var paymentStatus string
	queryPaymentStatus := `SELECT status FROM payments WHERE id = $1`
	err := pool.QueryRow(context.Background(), queryPaymentStatus, invoice.ID).Scan(&paymentStatus)
	if err != nil {
		return errors.New("payment not found")
	}

	if paymentStatus != "completed" {
		return errors.New("cannot create an invoice for a payment that is not completed")
	}

	query := `INSERT INTO invoices (id, amount, issueDate, dueDate, userId) VALUES ($1, $2, $3, $4, $5)`
	_, err = pool.Exec(context.Background(), query, invoice.ID, invoice.Amount, invoice.IssueDate, invoice.DueDate, invoice.UserID)
	if err != nil {
		return err
	}

	return nil
}

func GetInvoiceByID(invoiceID string, pool *pgxpool.Pool) (*models.Invoice, error) {
	query := `SELECT id, amount, issueDate, dueDate, userId FROM invoices WHERE id = $1`
	row := pool.QueryRow(context.Background(), query, invoiceID)

	var invoice models.Invoice
	err := row.Scan(&invoice.ID, &invoice.Amount, &invoice.IssueDate, &invoice.DueDate, &invoice.UserID)
	if err != nil {
		return nil, errors.New("invoice not found")
	}

	return &invoice, nil
}

func GetInvoicesByUserID(userID string, pool *pgxpool.Pool) ([]models.Invoice, error) {
	query := `SELECT id, amount, issueDate, dueDate, userId FROM invoices WHERE userId = $1`
	rows, err := pool.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []models.Invoice
	for rows.Next() {
		var invoice models.Invoice
		if err := rows.Scan(&invoice.ID, &invoice.Amount, &invoice.IssueDate, &invoice.DueDate, &invoice.UserID); err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

func UpdateInvoice(invoice *models.Invoice, pool *pgxpool.Pool) error {
	query := `UPDATE invoices SET amount = $2, issueDate = $3, dueDate = $4, userId = $5 WHERE id = $1`
	cmdTag, err := pool.Exec(context.Background(), query, invoice.ID, invoice.Amount, invoice.IssueDate, invoice.DueDate, invoice.UserID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("no rows were updated, invoice not found")
	}

	return nil
}

func DeleteInvoice(invoiceID string, pool *pgxpool.Pool) error {
	query := `DELETE FROM invoices WHERE id = $1`
	cmdTag, err := pool.Exec(context.Background(), query, invoiceID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("no rows were deleted, invoice not found")
	}

	return nil
}

func GenerateInvoiceForPayment(paymentID string, userID string, pool *pgxpool.Pool) (*models.Invoice, error) {
	// Fetch payment details
	var payment models.Payment
	paymentQuery := `SELECT id, amount, date, status FROM payments WHERE id = $1`
	err := pool.QueryRow(context.Background(), paymentQuery, paymentID).Scan(
		&payment.ID,
		&payment.Amount,
		&payment.Date,
		&payment.Status,
	)
	if err != nil {
		return nil, errors.New("payment not found")
	}

	// Check if the payment status is valid for invoice generation
	if payment.Status != "completed" {
		return nil, errors.New("cannot generate an invoice for a payment that is not completed")
	}

	// Parse payment date string into time.Time object
	paymentDate, err := time.Parse("2006-01-02", payment.Date) // Assumes `payment.Date` is stored as "YYYY-MM-DD" in the database
	if err != nil {
		return nil, errors.New("invalid payment date format")
	}

	// Create invoice details
	invoice := &models.Invoice{
		ID:        utils.GenerateUniqueID(),
		Amount:    payment.Amount,
		IssueDate: paymentDate.Format("2006-01-02"),                  // Format as "YYYY-MM-DD"
		DueDate:   paymentDate.AddDate(0, 1, 0).Format("2006-01-02"), // Add 1 month to the payment date
		UserID:    userID,
	}

	// Save the invoice to the database
	err = AddInvoice(invoice, pool)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

// Admin use
func GetAllInvoices(pool *pgxpool.Pool) ([]models.Invoice, error) {

	query := `SELECT id, amount, issueDate, dueDate, userId FROM invoices`
	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Slice to hold the invoices
	var invoices []models.Invoice
	for rows.Next() {
		var invoice models.Invoice
		// Scan each row into the invoice struct
		if err := rows.Scan(&invoice.ID, &invoice.Amount, &invoice.IssueDate, &invoice.DueDate, &invoice.UserID); err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}

	// Check for any errors during row iteration
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return invoices, nil
}
