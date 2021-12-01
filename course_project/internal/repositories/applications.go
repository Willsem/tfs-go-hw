package repositories

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/willsem/tfs-go-hw/course_project/internal/domain"
)

type PostgresqlApplicationsReposity struct {
	pool IPgxPool
}

func NewApplicaitionsRepository(pool IPgxPool) *PostgresqlApplicationsReposity {
	return &PostgresqlApplicationsReposity{
		pool: pool,
	}
}

const addApplicationQuery = `insert into applications (ticker, cost, size, type) values ($1, $2, $3, $4) returning (ticker)`

func (repository *PostgresqlApplicationsReposity) Add(application domain.Application) error {
	var row string
	err := repository.pool.QueryRow(context.Background(), addApplicationQuery,
		application.Ticker, application.Cost, application.Size, application.Type).Scan(&row)
	if err != nil {
		return err
	}

	return nil
}

const selectApplicationsQuery = `select id, ticker, cost, size, created_at, type from applications`

func (repository *PostgresqlApplicationsReposity) GetAll() ([]domain.Application, error) {
	rows, err := repository.pool.Query(context.Background(), selectApplicationsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return repository.scanApplications(rows)
}

const selectApplicationsByTicketQuery = selectApplicationsQuery + `where ticker=$1`

func (repository *PostgresqlApplicationsReposity) GetByTicker(ticker string) ([]domain.Application, error) {
	rows, err := repository.pool.Query(context.Background(), selectApplicationsByTicketQuery, ticker)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return repository.scanApplications(rows)
}

func (repository *PostgresqlApplicationsReposity) scanApplications(rows pgx.Rows) ([]domain.Application, error) {
	applications := make([]domain.Application, 0)

	for rows.Next() {
		var app domain.Application
		err := rows.Scan(&app.Id, &app.Ticker, &app.Cost, &app.Size, &app.CreatedAt, &app.Type)
		if err != nil {
			return nil, err
		}

		applications = append(applications, app)
	}

	return applications, nil
}
