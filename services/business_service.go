package services

import (
	"context"
	"errors"
	"github.com/Bradkibs/MONOS-challenge/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllBusinesses(pool *pgxpool.Pool) ([]models.Business, error) {
	rows, err := pool.Query(context.Background(), "SELECT id, vendor_id, name, description, deleted_at FROM businesses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var businesses []models.Business
	for rows.Next() {
		var business models.Business
		if err := rows.Scan(&business.ID, &business.VendorID, &business.Name, &business.Description, &business.DeletedAt); err != nil {
			return nil, err
		}
		businesses = append(businesses, business)
	}

	return businesses, nil
}

func CreateBusiness(business *models.Business, pool *pgxpool.Pool) error {
	existingQuery := `SELECT COUNT(*) FROM businesses WHERE name = $1 AND vendor_id = $2`
	var count int
	err := pool.QueryRow(context.Background(), existingQuery, business.Name, business.VendorID).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("business with the same name already exists for this vendor")
	}

	query := `INSERT INTO businesses (id, vendor_id, name, description, deleted_at) VALUES ($1, $2, $3, $4, $5)`
	_, err = pool.Exec(context.Background(), query, business.ID, business.VendorID, business.Name, business.Description, business.DeletedAt)
	return err
}

func GetBusinessByID(businessID uuid.UUID, pool *pgxpool.Pool) (*models.Business, error) {
	query := `SELECT id, vendor_id, name, description, deleted_at FROM businesses WHERE id = $1`
	row := pool.QueryRow(context.Background(), query, businessID)

	var business models.Business
	if err := row.Scan(&business.ID, &business.VendorID, &business.Name, &business.Description, &business.DeletedAt); err != nil {
		return nil, err
	}

	return &business, nil
}

func UpdateBusiness(business *models.Business, pool *pgxpool.Pool) error {
	query := `UPDATE businesses SET vendor_id = $2, name = $3, description = $4, deleted_at = $5 WHERE id = $1`
	cmdTag, err := pool.Exec(context.Background(), query, business.ID, business.VendorID, business.Name, business.Description, business.DeletedAt)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("no rows were updated, business not found")
	}

	return nil
}

func DeleteBusiness(businessID uuid.UUID, pool *pgxpool.Pool) error {
	query := `DELETE FROM businesses WHERE id = $1`
	cmdTag, err := pool.Exec(context.Background(), query, businessID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("no rows were deleted, business not found")
	}

	return nil
}

func GetBusinessesByVendorID(vendorID uuid.UUID, pool *pgxpool.Pool) ([]models.Business, error) {
	query := `SELECT id, vendor_id, name, description, deleted_at FROM businesses WHERE vendor_id = $1`
	rows, err := pool.Query(context.Background(), query, vendorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var businesses []models.Business
	for rows.Next() {
		var business models.Business
		if err := rows.Scan(&business.ID, &business.VendorID, &business.Name, &business.Description, &business.DeletedAt); err != nil {
			return nil, err
		}
		businesses = append(businesses, business)
	}

	return businesses, nil
}
