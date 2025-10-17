package models

import "time"

type Transaction struct {
	ID string `json:"id"`
	Quantity int `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
	Date time.Time `json:"time"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}