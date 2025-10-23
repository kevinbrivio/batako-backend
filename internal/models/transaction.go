package models

import "time"

type Transaction struct {
	ID string `json:"id"`
	Quantity int `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
	PurchaseDate time.Time `json:"purchase_date"`
	Customer string `json:"customer"`
	Address string `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}