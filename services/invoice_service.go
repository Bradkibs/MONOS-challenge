package services

import (
	"context"
	"errors"
	"time"

	"github.com/Bradkibs/MONOS-challenge/models"
	"github.com/Bradkibs/MONOS-challenge/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AddInvoice(invoice *models.Invoice, pool *pgxpool.Pool) error {
	var paymentStatus string
	queryPaymentStatus := `SELECT status FROM payments WHERE id = $1`
	err := pool.QueryRow(context.Background(), queryPaymentStatus, invoice.PaymentID).Scan(&paymentStatus)
	if err != nil {
		return errors.New("payment not found")
	}

	if paymentStatus != "completed" {
		return errors.New("cannot create an invoice for a payment that is not completed")
	}

	query := `INSERT INTO invoices (id, payment_id, issue_date, due_date, status) VALUES ($1, $2, $3, $4, $5)`
	_, err = pool.Exec(context.Background(), query, invoice.ID, invoice.PaymentID, invoice.IssueDate, invoice.DueDate, invoice.Status)
	return err
}

func GetInvoiceByID(invoiceID uuid.UUID, pool *pgxpool.Pool) (*models.Invoice, error) {
	query := `SELECT id, payment_id, issue_date, due_date, status, deleted_at FROM invoices WHERE id = $1`
	row := pool.QueryRow(context.Background(), query, invoiceID)

	var invoice models.Invoice
	err := row.Scan(&invoice.ID, &invoice.PaymentID, &invoice.IssueDate, &invoice.DueDate, &invoice.Status, &invoice.DeletedAt)
	if err != nil {
		return nil, errors.New("invoice not found")
	}

	return &invoice, nil
}

func UpdateInvoice(invoice *models.Invoice, pool *pgxpool.Pool) error {
	query := `UPDATE invoices SET payment_id = $2, issue_date = $3, due_date = $4, status = $5 WHERE id = $1`
	cmdTag, err := pool.Exec(context.Background(), query, invoice.ID, invoice.PaymentID, invoice.IssueDate, invoice.DueDate, invoice.Status)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("no rows were updated, invoice not found")
	}

	return nil
}

func DeleteInvoice(invoiceID uuid.UUID, pool *pgxpool.Pool) error {
	query := `UPDATE invoices SET deleted_at = $2 WHERE id = $1`
	cmdTag, err := pool.Exec(context.Background(), query, invoiceID, time.Now())
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("no rows were deleted, invoice not found")
	}

	return nil
}

func GenerateInvoiceForPayment(paymentID, userID uuid.UUID, pool *pgxpool.Pool) (*models.Invoice, error) {
	var payment models.Payment
	paymentQuery := `SELECT id, amount, date, status FROM payments WHERE id = $1`
	err := pool.QueryRow(context.Background(), paymentQuery, paymentID).Scan(&payment.ID, &payment.Amount, &payment.Date, &payment.Status)
	if err != nil {
		return nil, errors.New("payment not found")
	}

	if payment.Status != "completed" {
		return nil, errors.New("cannot generate an invoice for a payment that is not completed")
	}

	invoice := &models.Invoice{
		ID:        utils.GenerateUniqueID(),
		PaymentID: payment.ID,
		IssueDate: payment.Date,
		DueDate:   payment.Date.AddDate(0, 1, 0),
		Status:    "issued",
	}

	err = AddInvoice(invoice, pool)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func GetAllInvoices(pool *pgxpool.Pool) ([]models.Invoice, error) {
	query := `SELECT id, payment_id, issue_date, due_date, status, deleted_at FROM invoices`
	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []models.Invoice
	for rows.Next() {
		var invoice models.Invoice
		if err := rows.Scan(&invoice.ID, &invoice.PaymentID, &invoice.IssueDate, &invoice.DueDate, &invoice.Status, &invoice.DeletedAt); err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}

	return invoices, nil
}
