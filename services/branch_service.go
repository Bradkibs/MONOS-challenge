package services

import (
	"context"
	"errors"
	"github.com/Bradkibs/MONOS-challenge/models"
	"github.com/Bradkibs/MONOS-challenge/utils"
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
func UpdateBranchesForSubscription(subscriptionID string, branchChange int, branchNames []string, pool *pgxpool.Pool) error {
	// Fetch the business ID for the subscription
	var businessID string
	businessQuery := `SELECT businessId FROM subscriptions WHERE id = $1`
	err := pool.QueryRow(context.Background(), businessQuery, subscriptionID).Scan(&businessID)
	if err != nil {
		return errors.New("subscription not found or invalid")
	}

	// Fetch the current branch count for the business
	var branchCount int
	countQuery := `SELECT COUNT(*) FROM branches WHERE businessId = $1`
	err = pool.QueryRow(context.Background(), countQuery, businessID).Scan(&branchCount)
	if err != nil {
		return err
	}

	// Calculate the new branch count
	newBranchCount := branchCount + branchChange
	if newBranchCount < 1 {
		return errors.New("cannot remove more branches than currently exist")
	}

	// Adding branches
	if branchChange > 0 {
		if len(branchNames) < branchChange {
			return errors.New("not enough branch names provided for the number of branches to add")
		}

		for _, branchName := range branchNames[:branchChange] {
			newBranchID := utils.GenerateUniqueID()
			insertQuery := `INSERT INTO branches (id, businessId, location) VALUES ($1, $2, $3)`
			_, err := pool.Exec(context.Background(), insertQuery, newBranchID, businessID, branchName)
			if err != nil {
				return err
			}
		}
	}

	// Removing branches
	if branchChange < 0 {
		if len(branchNames) < -branchChange {
			return errors.New("not enough branch names provided for the number of branches to remove")
		}

		for _, branchName := range branchNames[:(-branchChange)] {
			deleteQuery := `DELETE FROM branches WHERE businessId = $1 AND location = $2`
			_, err := pool.Exec(context.Background(), deleteQuery, businessID, branchName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
