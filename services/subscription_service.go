package services

import (
	"context"
	"errors"
	"time"

	"github.com/Bradkibs/MONOS-challenge/models"
	"github.com/jackc/pgx/v5/pgxpool"
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
	totalCost := basePrice + float64(branchCount)*1.0
	return totalCost, nil
}

func CancelSubscription(subscriptionID string, pool *pgxpool.Pool) error {
	query := `SELECT startDate, status FROM subscriptions WHERE id = $1 AND deleted_at IS NULL`
	var startDate time.Time
	var status string

	err := pool.QueryRow(context.Background(), query, subscriptionID).Scan(&startDate, &status)
	if err != nil {
		return err
	}

	if status != "active" {
		return errors.New("subscription is not active, cancellation not possible")
	}

	if time.Since(startDate) > 7*24*time.Hour {
		return errors.New("refund not allowed after 1 week of subscription start")
	}

	updateQuery := `UPDATE subscriptions SET status = 'canceled', deleted_at = NOW() WHERE id = $1`
	_, err = pool.Exec(context.Background(), updateQuery, subscriptionID)
	return err
}

func DowngradeSubscription(subscriptionID, newTier string, pool *pgxpool.Pool) error {
	query := `SELECT s.businessId, COUNT(p.id) FROM subscriptions s LEFT JOIN products p ON s.businessId = p.businessId WHERE s.id = $1 AND s.deleted_at IS NULL GROUP BY s.businessId`
	var businessID string
	var productCount int

	err := pool.QueryRow(context.Background(), query, subscriptionID).Scan(&businessID, &productCount)
	if err != nil {
		return err
	}

	if (newTier == "Starter" && productCount > 10) || (newTier == "Pro" && productCount > 100) {
		return errors.New("reduce product count before downgrading")
	}

	updateQuery := `UPDATE subscriptions SET tier = $2 WHERE id = $1 AND deleted_at IS NULL`
	_, err = pool.Exec(context.Background(), updateQuery, subscriptionID, newTier)
	return err
}

func HandleSubscriptionOverlap(currentSubscriptionID string, newSubscription *models.Subscription, pool *pgxpool.Pool) error {
	query := `SELECT endDate, status FROM subscriptions WHERE id = $1 AND deleted_at IS NULL`
	var currentEndDate time.Time
	var currentStatus string

	err := pool.QueryRow(context.Background(), query, currentSubscriptionID).Scan(&currentEndDate, &currentStatus)
	if err != nil {
		return err
	}

	if currentStatus != "active" {
		return errors.New("current subscription is not active")
	}

	if newSubscription.StartDate.Before(currentEndDate) {
		extendedEndDate := currentEndDate.Add(newSubscription.EndDate.Sub(newSubscription.StartDate))
		updateQuery := `UPDATE subscriptions SET endDate = $2 WHERE id = $1 AND deleted_at IS NULL`
		_, err = pool.Exec(context.Background(), updateQuery, currentSubscriptionID, extendedEndDate)
		return err
	}

	return errors.New("no overlap detected")
}

func CreateSubscription(subscription *models.Subscription, pool *pgxpool.Pool) error {
	query := `INSERT INTO subscriptions (id, businessId, tier, startDate, endDate, status) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := pool.Exec(context.Background(), query, subscription.ID, subscription.BusinessID, subscription.Tier, subscription.StartDate, subscription.EndDate, subscription.Status)
	return err
}

func GetSubscription(subscriptionID string, pool *pgxpool.Pool) (*models.Subscription, error) {
	query := `SELECT id, businessId, tier, startDate, endDate, status FROM subscriptions WHERE id = $1 AND deleted_at IS NULL`
	var subscription models.Subscription
	err := pool.QueryRow(context.Background(), query, subscriptionID).Scan(
		&subscription.ID,
		&subscription.BusinessID,
		&subscription.Tier,
		&subscription.StartDate,
		&subscription.EndDate,
		&subscription.Status,
	)
	if err != nil {
		return nil, errors.New("subscription not found")
	}
	return &subscription, nil
}

func UpdateSubscription(subscription *models.Subscription, pool *pgxpool.Pool) error {
	query := `UPDATE subscriptions SET tier = $2, startDate = $3, endDate = $4, status = $5 WHERE id = $1 AND deleted_at IS NULL`
	cmdTag, err := pool.Exec(context.Background(), query, subscription.ID, subscription.Tier, subscription.StartDate, subscription.EndDate, subscription.Status)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("no rows updated")
	}
	return nil
}

func DeleteSubscription(subscriptionID string, pool *pgxpool.Pool) error {
	query := `UPDATE subscriptions SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	cmdTag, err := pool.Exec(context.Background(), query, subscriptionID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("no rows deleted")
	}
	return nil
}
