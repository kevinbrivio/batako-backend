package models

import (
	"time"
)

type CementType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CementStock struct {
	ID           string     `json:"id,omitempty" db:"id"`
	CementType   CementType `json:"cement_type" db:"-"`
	Quantity     int        `json:"quantity" db:"quantity"`
	PricePerBag  float64    `json:"price_per_bag" db:"price_per_bag"`
	TotalPrice   float64    `json:"total_price" db:"-"`
	PurchaseDate time.Time  `json:"purchase_date" db:"purchase_date"`
	CreatedAt    time.Time  `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at,omitempty" db:"updated_at"`
}

type CreateCementStockRequest struct {
	CementTypeName string    `json:"cement_type_name"`
	Quantity       int       `json:"quantity"`
	PricePerBag    float64   `json:"price_per_bag"`
	PurchaseDate   time.Time `json:"purchase_date"`
}

type UpdateCementStockRequest struct {
	ID             string    `json:"id"`
	CementTypeName string    `json:"cement_type_name"`
	Quantity       int       `json:"quantity"`
	PricePerBag    float64   `json:"price_per_bag"`
	PurchaseDate   time.Time `json:"purchase_date"`
}
