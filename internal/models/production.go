package models

import "time"

type Production struct {
	ID string `json:"id"`
	Quantity int `json:"quantity"`
	CementUsed float64 `json:"cement_used"`
	SandUsed float64 `json:"sand_used"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}