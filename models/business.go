package models

type Business struct {
	ID             string `json:"id"`
	VendorID       string `json:"vendor_id"`
	Name           string `json:"name"`
	SubscriptionID string `json:"subscription_id"`
}
