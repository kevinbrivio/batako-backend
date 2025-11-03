package store

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/kevinbrivio/batako-backend/internal/models"
	"github.com/kevinbrivio/batako-backend/internal/utils"
)

type SalaryStore struct {
	db *sql.DB
	prodStore *ProductionStore
}

func (s *SalaryStore) GetWeekly(ctx context.Context, productionDate time.Time) (float64, int, error) {
	targetDate := time.Date(productionDate.Year(), productionDate.Month(), productionDate.Day(), 0, 0, 0, 0, productionDate.Location())
	
	query := `
		SELECT salary, total_production
		FROM employee_salary
		WHERE $1 BETWEEN start_date AND end_date
	`
	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()

	var salary float64
	var totalProduction int 

	err := s.db.QueryRowContext(ctx, query, targetDate).Scan(&salary, &totalProduction)
	if err != nil {
		return 0, 0, nil
	}

	log.Printf("Salary: %.2f", salary)

	return salary, totalProduction, nil
}

func (s *SalaryStore) GetMonthly(ctx context.Context, monthOffset int) ([]models.EmployeeSalary, int, error) {
	today := time.Now()

	start, end := utils.GetMonthRange(today, monthOffset)
	
	query := `
		SELECT id, salary, total_production, start_date, end_date, count(*) over() as total_count
		FROM employee_salary
		WHERE start_date <= $2 AND end_date >= $1
	`
	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	salaries := []models.EmployeeSalary{}
	var totalCount int

	for rows.Next() {
		var sal models.EmployeeSalary
		if err := rows.Scan(
			&sal.ID,
			&sal.Salary,
			&sal.TotalProduction,
			&sal.StartDate,
			&sal.EndDate,
			&totalCount,
		); err != nil {
			return salaries, 0, err
		}
		salaries = append(salaries, sal)
	}

	return salaries, totalCount, nil
}

func (s *SalaryStore) AddSalary(ctx context.Context, es *models.EmployeeSalary) error {
	es.ID = uuid.New().String()
	
	query := `
		INSERT INTO employee_salary (id, start_date, end_date, total_production, salary)
		VALUES ($1, $2, $3, $4, $5) RETURNING created_at, updated_at	
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx, query,
		es.ID,
		es.StartDate,
		es.EndDate,
		es.TotalProduction,
		es.Salary,
	).Scan(
		&es.CreatedAt,
		&es.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *SalaryStore) GenerateWeeklySalary(ctx context.Context) error {
	today := time.Now()
	weekday := int(today.Weekday())

	if weekday == 0 { weekday = 7 }

    // Calculate week range
    startDate := today.AddDate(0, 0, -weekday + 1)
    endDate := startDate.AddDate(0, 0, 6)
	
    totalProduction, err := s.prodStore.GetTotalProduction(ctx, startDate, endDate)
    if err != nil {
        return err
    }

	const salaryBase = 450
    salary := float64(totalProduction) * salaryBase

	employeeSalary := models.EmployeeSalary{
		StartDate: startDate,
		EndDate: endDate,
		TotalProduction: totalProduction,
		Salary: salary,
	}

    // Insert into employee_salary
    return s.AddSalary(ctx, &employeeSalary)
}

func(s *SalaryStore) StartSchedulers(ctx context.Context) {
    scheduler, _ := gocron.NewScheduler(gocron.WithLocation(time.Local))

    // Run every Saturday at 23:59 (11:59 PM)
    _, err := scheduler.NewJob(
        gocron.WeeklyJob(
			1, // One week a time
			gocron.NewWeekdays(time.Wednesday),
			gocron.NewAtTimes(gocron.NewAtTime(17, 16, 0)),
		),
		gocron.NewTask(func() {
			if err := s.GenerateWeeklySalary(ctx); err != nil {
				log.Printf("Failed to generate weekly salary")
			} else {
				log.Println("Weekly salary successfully generated")
			}
		}),
    )
    if err != nil {
        log.Fatalf("Failed to create job: %v", err)
    }

    scheduler.Start()
}
