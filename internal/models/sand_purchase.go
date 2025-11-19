package models

import "time"

type SandPurchase struct {
	ID string `json:"id"`
	SandType string `json:"sand_type"`
	Quantity int `json:"quantity"`
	PricePerTruck float64 `json:"price_per_truck"`
	PurchaseDate time.Time `json:"purchase_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}