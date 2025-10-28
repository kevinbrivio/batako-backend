package models
import "time"

type EmployeeSalary struct {
	ID string `json:"id"`
	StartDate time.Time `json:"start_date"`
	EndDate time.Time `json:"end_date"`
	TotalProduction int `json:"total_production"`
	Salary float64 `json:"salary"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
