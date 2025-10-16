package models

import "time"

type Production struct {
	ID string `json:"id"`
	Quantity int `json:"quantity"`
	CementUsed float64 `json:"cement_used"`
	Date float64 `json:"date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}