package models

import "time"

type Mutation struct {
	ID        uint      `json:"id"`
	AccountID uint      `json:"account_id"`
	Nominal   float64   `json:"nominal"`
	Type      string    `json:"type"`
	Reference string    `json:"reference"`
	CreatedAt time.Time `json:"created_at"`
}
