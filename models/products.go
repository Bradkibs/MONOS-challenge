package models

type Product struct {
	ID         string  `json:"id"`
	BusinessID string  `json:"businessId"`
	Name       string  `json:"name"`
	Details    string  `json:"details"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
}
