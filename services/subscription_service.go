package services

import (
	"context"
	"errors"
	"github.com/Bradkibs/MONOS-challenge/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func CalculateSubscriptionCost(subscriptionTier string, branchCount int) (float64, error) {
	var basePrice float64

	switch subscriptionTier {
	case "Starter":
		basePrice = 1.0
	case "Pro":
		basePrice = 3.0
	case "Enterprise":
		basePrice = 5.0
	default:
		return 0.0, errors.New("invalid subscription tier")
	}

	// Calculate the dynamic pricing for additional branches
	totalCost := basePrice + float64(branchCount)*1.0
	return totalCost, nil
}

func CancelSubscription(subscriptionID string, pool *pgxpool.Pool) error {
	query := `SELECT startDate, status FROM subscriptions WHERE id = $1`
	var startDate time.Time
	var status string

	// Fetch subscription details
	err := pool.QueryRow(context.Background(), query, subscriptionID).Scan(&startDate, &status)
	if err != nil {
		return err
	}

	// Ensure the subscription is active
	if status != "active" {
		return errors.New("subscription is not active, cancellation not possible")
	}

	// Check if cancellation is within 1 week of the start date
	if time.Since(startDate) > 7*24*time.Hour {
		return errors.New("refund not allowed after 1 week of subscription start")
	}

	// Deactivate subscription and associated listings
	updateSubscription := `UPDATE subscriptions SET status = 'canceled' WHERE id = $1`
	_, err = pool.Exec(context.Background(), updateSubscription, subscriptionID)
	if err != nil {
		return err
	}

	// No refund is issued, but listings are deactivated
	return nil
}

func DowngradeSubscription(subscriptionID string, newTier string, pool *pgxpool.Pool) error {
	// Fetch business ID and product count for the subscription
	query := `SELECT s.businessId, COUNT(p.id) FROM subscriptions s LEFT JOIN products p ON s.businessId = p.businessId WHERE s.id = $1 GROUP BY s.businessId`
	var businessID string
	var productCount int

	err := pool.QueryRow(context.Background(), query, subscriptionID).Scan(&businessID, &productCount)
	if err != nil {
		return err
	}

	// Validate product count for downgrade
	if newTier == "Starter" && productCount > 10 {
		return errors.New("reduce product count to 10 or fewer before downgrading to Starter")
	}
	// Validate product count for downgrade
	if newTier == "Pro" && productCount > 100 {
		return errors.New("reduce product count to 100 or fewer before downgrading to Pro")
	}

	// Perform the downgrade
	updateQuery := `UPDATE subscriptions SET tier = $2 WHERE id = $1`
	_, err = pool.Exec(context.Background(), updateQuery, subscriptionID, newTier)
	if err != nil {
		return err
	}

	return nil
}

func HandleSubscriptionOverlap(currentSubscriptionID string, newSubscription *models.Subscription, pool *pgxpool.Pool) error {
	query := `SELECT endDate, status FROM subscriptions WHERE id = $1`
	var currentEndDate time.Time
	var currentStatus string

	// Fetch current subscription details
	err := pool.QueryRow(context.Background(), query, currentSubscriptionID).Scan(&currentEndDate, &currentStatus)
	if err != nil {
		return err
	}

	// Ensure the current subscription is active
	if currentStatus != "active" {
		return errors.New("current subscription is not active, cannot merge")
	}

	// Parse new subscription dates as time.Time
	newStartDate, err := time.Parse("2016-01-02", newSubscription.StartDate)
	if err != nil {
		return errors.New("invalid new subscription start date format")
	}
	newEndDate, err := time.Parse("2016-01-02", newSubscription.EndDate)
	if err != nil {
		return errors.New("invalid new subscription end date format")
	}

	// Check for overlap and merge subscriptions
	if newStartDate.Before(currentEndDate) {
		// Extend the current subscription's end date
		extendedEndDate := currentEndDate.Add(newEndDate.Sub(newStartDate))
		updateQuery := `UPDATE subscriptions SET endDate = $2 WHERE id = $1`
		_, err = pool.Exec(context.Background(), updateQuery, currentSubscriptionID, extendedEndDate)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("no overlap detected or invalid subscription")
}
