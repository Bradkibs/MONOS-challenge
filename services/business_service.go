package services

import (
	"context"
	"errors"
	"github.com/Bradkibs/MONOS-challenge/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllBusinesses(pool *pgxpool.Pool) ([]models.Business, error) {
	rows, err := pool.Query(context.Background(), "SELECT id, vendor_id, name, subscription_id FROM businesses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var businesses []models.Business
	for rows.Next() {
		var business models.Business
		if err := rows.Scan(&business.ID, &business.VendorID, &business.Name, &business.SubscriptionID); err != nil {
			return nil, err
		}
		businesses = append(businesses, business)
	}

	return businesses, nil
}

func CreateBusiness(business *models.Business, pool *pgxpool.Pool) error {
	// Check if a business with the same name and vendor_id already exists
	existingQuery := `SELECT COUNT(*) FROM businesses WHERE name = $1 AND vendor_id = $2`
	var count int
	err := pool.QueryRow(context.Background(), existingQuery, business.Name, business.VendorID).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("business with the same name already exists for this vendor")
	}

	query := `INSERT INTO businesses (id, vendor_id, name, subscription_id) VALUES ($1, $2, $3, $4)`
	_, err = pool.Exec(context.Background(), query, business.ID, business.VendorID, business.Name, business.SubscriptionID)
	if err != nil {
		return err
	}

	return nil
}

func GetBusinessByID(businessID string, pool *pgxpool.Pool) (*models.Business, error) {
	query := `SELECT id, vendor_id, name, subscription_id FROM businesses WHERE id = $1`
	row := pool.QueryRow(context.Background(), query, businessID)

	var business models.Business
	if err := row.Scan(&business.ID, &business.VendorID, &business.Name, &business.SubscriptionID); err != nil {
		return nil, err
	}

	return &business, nil
}

func GetBusinessByName(name string, pool *pgxpool.Pool) (*models.Business, error) {
	query := `SELECT id, vendor_id, name, subscription_id FROM businesses WHERE name = $1`
	row := pool.QueryRow(context.Background(), query, name)

	var business models.Business
	if err := row.Scan(&business.ID, &business.VendorID, &business.Name, &business.SubscriptionID); err != nil {
		return nil, err
	}

	return &business, nil
}

func GetBusinessByVendorId(vendorId string, pool *pgxpool.Pool) (*models.Business, error) {
	query := `SELECT id, vendor_id, name, subscription_id FROM businesses WHERE vendor_id = $1`
	row := pool.QueryRow(context.Background(), query, vendorId)

	var business models.Business
	if err := row.Scan(&business.ID, &business.VendorID, &business.Name, &business.SubscriptionID); err != nil {
		return nil, err
	}

	return &business, nil
}

func GetBusinessBySubscriptionId(subscriptionId string, pool *pgxpool.Pool) (*models.Business, error) {
	query := `SELECT id, vendor_id, name, subscription_id FROM businesses WHERE subscription_id = $1`
	row := pool.QueryRow(context.Background(), query, subscriptionId)

	var business models.Business
	if err := row.Scan(&business.ID, &business.VendorID, &business.Name, &business.SubscriptionID); err != nil {
		return nil, err
	}

	return &business, nil
}
func UpdateBusiness(business *models.Business, pool *pgxpool.Pool) error {
	query := `UPDATE businesses SET vendor_id = $2, name = $3, subscription_id = $4 WHERE id = $1`
	cmdTag, err := pool.Exec(context.Background(), query, business.ID, business.VendorID, business.Name, business.SubscriptionID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("no rows were updated, business not found")
	}

	return nil
}

// Delete a business by ID
func DeleteBusiness(businessID string, pool *pgxpool.Pool) error {
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

func GetBusinessesByVendorID(vendorID string, pool *pgxpool.Pool) ([]models.Business, error) {
	query := `SELECT id, vendor_id, name, subscription_id FROM businesses WHERE vendor_id = $1`
	rows, err := pool.Query(context.Background(), query, vendorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var businesses []models.Business
	for rows.Next() {
		var business models.Business
		if err := rows.Scan(&business.ID, &business.VendorID, &business.Name, &business.SubscriptionID); err != nil {
			return nil, err
		}
		businesses = append(businesses, business)
	}

	return businesses, nil
}
