package models

type DashboardResponse struct {
	Status   string           `json:"status"`
	Message  string           `json:"message"`
	Data     Dashboard    		`json:"data"`
}

type Dashboard struct {
	CementSummary CementSummary `json:"cement_summary"`
	SandSummary SandSummary `json:"sand_summary"`
	SalarySummary SalarySummary `json:"salary_summary"`
	ProductionSummary ProductionSummary `json:"production_summary"`
	TransactionSummary TransactionSummary `json:"transaction_summary"`
}

type CementSummary struct {
	TotalStock int `json:"total_stock"`
	TotalQuantity int `json:"total_quantity"`
	TotalPrice float64 `json:"total_price"`
}

type SandSummary struct {
	TotalPurchase int `json:"total_purchase"`
	TotalQuantity int `json:"total_quantity"`
	TotalPrice float64 `json:"total_price"`
}

type ProductionSummary struct {
	TotalProduction int `json:"total_production"`
}

type TransactionSummary struct {
	TotalTransaction int `json:"total_transaction"`
	TotalIncome float64 `json:"total_income"`
}

type SalarySummary struct {
	TotalSalary float64 `json:"total_salary"`
}
