package services

import (
	"context"
	"errors"
	"github.com/Bradkibs/MONOS-challenge/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AddBranch(branch *models.Branch, pool *pgxpool.Pool) error {
	query := `INSERT INTO branches (id, business_id, location) VALUES ($1, $2, $3)`
	_, err := pool.Exec(context.Background(), query, branch.ID, branch.BusinessID, branch.Location)
	if err != nil {
		return err
	}
	return nil
}

func GetBranchesByBusinessID(businessID string, pool *pgxpool.Pool) ([]models.Branch, error) {
	query := `SELECT id, business_id, location FROM branches WHERE business_id = $1`
	rows, err := pool.Query(context.Background(), query, businessID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var branches []models.Branch
	for rows.Next() {
		var branch models.Branch
		if err := rows.Scan(&branch.ID, &branch.BusinessID, &branch.Location); err != nil {
			return nil, err
		}
		branches = append(branches, branch)
	}

	return branches, nil
}

func UpdateBranch(branch *models.Branch, pool *pgxpool.Pool) error {
	query := `UPDATE branches SET location = $2 WHERE id = $1 AND business_id = $3`
	cmdTag, err := pool.Exec(context.Background(), query, branch.ID, branch.Location, branch.BusinessID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("no rows were updated, branch not found")
	}

	return nil
}

func DeleteBranch(branchID string, businessID string, pool *pgxpool.Pool) error {
	query := `DELETE FROM branches WHERE id = $1 AND business_id = $2`
	cmdTag, err := pool.Exec(context.Background(), query, branchID, businessID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("no rows were deleted, branch not found")
	}

	return nil
}
func UpdateBranchesForSubscription(subscriptionID string, branchChange int, pool *pgxpool.Pool) error {
	// Fetch branch count for the business associated with the subscription
	query := `
		SELECT COUNT(*) 
		FROM branches 
		WHERE businessId = (SELECT businessId FROM subscriptions WHERE id = $1)`
	var branchCount int

	err := pool.QueryRow(context.Background(), query, subscriptionID).Scan(&branchCount)
	if err != nil {
		return err
	}

	// Calculate the new branch count
	newBranchCount := branchCount + branchChange
	if newBranchCount < 0 {
		return errors.New("cannot remove more branches than currently exist")
	}

	// Prorate charges for added/removed branches (billing adjustment logic would be implemented here)
	// Assuming charges are calculated externally

	return nil
}
