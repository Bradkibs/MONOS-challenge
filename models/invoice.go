package models

type Invoice struct {
	ID        string  `json:"id"`
	Amount    float64 `json:"amount"`
	IssueDate string  `json:"issue_date"`
	DueDate   string  `json:"due_date"`
	UserID    string  `json:"user_id"`
}
