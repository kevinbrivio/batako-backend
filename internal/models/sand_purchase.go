package models

import "time"

type SandType struct {
	ID     int        `json:"id"`
	Name   string     `json:"name"`
}

type SandPurchase struct {
	ID              string      `json:"id"                db:"id"`
	SandType				SandType	  `json:"sand_type"         db:"-"`
	Quantity        int         `json:"quantity"          db:"quantity"`
	PricePerTruck   float64     `json:"price_per_truck"   db:"price_per_truck"`
	TotalPrice   		float64     `json:"total_price" db:"-"`
	PurchaseDate    time.Time   `json:"purchase_date"     db:"purchase_date"`
	CreatedAt       time.Time   `json:"created_at"        db:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"        db:"updated_at"`
}

type CreateSandPurchaseRequest struct {
	SandTypeName    string      `json:"sand_type_name"`
	Quantity        int         `json:"quantity"`
	PricePerTruck   float64     `json:"price_per_truck"`
	PurchaseDate    time.Time   `json:"purchase_date"`
}

type UpdateSandPurchaseRequest struct {
	ID              string      `json:"id"`
	SandTypeName    string      `json:"sand_type_name"`
	Quantity        int         `json:"quantity"`
	PricePerTruck   float64     `json:"price_per_truck"`
	PurchaseDate    time.Time   `json:"purchase_date"`
}
