package models

type Subscription struct {
	ID        string  `json:"id"`
	Tier      string  `json:"tier"`
	Price     float64 `json:"price"`
	StartDate string  `json:"start_date"`
	EndDate   string  `json:"end_date"`
	Status    string  `json:"Status"`
}
