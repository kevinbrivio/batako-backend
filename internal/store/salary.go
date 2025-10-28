package store

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/kevinbrivio/batako-backend/internal/models"
)

type SalaryStore struct {
	db *sql.DB
	prodStore *ProductionStore
}

func (s *SalaryStore) GetWeeklySalary(ctx context.Context, productionDate time.Time) (float64, error) {
	today := time.Now()
	targetDate := int(productionDate.Weekday())
	startDate := today.AddDate(0, 0, -targetDate + 1)
	endDate := startDate.AddDate(0, 0, 6)
	
	query := `
		SELECT salary
		FROM employee_salary
		WHERE start_date = $2 AND end_date = $3
	`
	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()

	var salary float64
	err := s.db.QueryRowContext(ctx, query, startDate, endDate).Scan(&salary)
	if err != nil {
		return 0, nil
	}

	return salary, nil
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
			gocron.NewWeekdays(time.Saturday),
			gocron.NewAtTimes(gocron.NewAtTime(23, 59, 59)),
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
