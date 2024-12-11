package services

import (
	"context"
	"errors"

	"github.com/Bradkibs/MONOS-challenge/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AddProduct(product *models.Product, pool *pgxpool.Pool) error {
	// Check if the product already exists (by ID or BusinessID and Name)
	checkQuery := `SELECT COUNT(*) FROM products WHERE id = $1 OR (businessId = $2 AND name = $3 AND deleted_at IS NULL)`
	var count int
	err := pool.QueryRow(context.Background(), checkQuery, product.ID, product.BusinessID, product.Name).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("product already exists")
	}

	// Insert product into the database
	insertQuery := `INSERT INTO products (id, businessId, name, details, quantity, price) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = pool.Exec(context.Background(), insertQuery, product.ID, product.BusinessID, product.Name, product.Details, product.Quantity, product.Price)
	if err != nil {
		return err
	}

	return nil
}

func GetProductsByBusinessID(businessID string, pool *pgxpool.Pool) ([]models.Product, error) {
	query := `SELECT id, businessId, name, details, quantity, price FROM products WHERE businessId = $1 AND deleted_at IS NULL`
	rows, err := pool.Query(context.Background(), query, businessID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.BusinessID, &product.Name, &product.Details, &product.Quantity, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func UpdateProduct(product *models.Product, pool *pgxpool.Pool) error {
	query := `UPDATE products SET name = $2, details = $3, quantity = $4, price = $5 WHERE id = $1 AND businessId = $6 AND deleted_at IS NULL`
	cmdTag, err := pool.Exec(context.Background(), query, product.ID, product.Name, product.Details, product.Quantity, product.Price, product.BusinessID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("no rows were updated, product not found")
	}

	return nil
}

func DeleteProduct(productID, businessID string, pool *pgxpool.Pool) error {
	query := `UPDATE products SET deleted_at = NOW() WHERE id = $1 AND businessId = $2 AND deleted_at IS NULL`
	cmdTag, err := pool.Exec(context.Background(), query, productID, businessID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("no rows were deleted, product not found")
	}

	return nil
}
