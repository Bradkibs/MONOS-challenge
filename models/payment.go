package models

type Payment struct {
	ID             string  `json:"id"`
	SubscriptionID string  `json:"subscription_id"`
	Amount         float64 `json:"amount"`
	Date           string  `json:"date"`
	Status         string  `json:"status"`
}
