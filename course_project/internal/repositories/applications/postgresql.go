package applications

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/willsem/tfs-go-hw/course_project/internal/domain"
)

type PostgresqlApplicationsReposity struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *PostgresqlApplicationsReposity {
	return &PostgresqlApplicationsReposity{
		pool: pool,
	}
}

const addApplicationQuery = `insert into applications (ticker, cost) values ($1, $2) returning (ticker)`

func (repository *PostgresqlApplicationsReposity) Add(application domain.Application) error {
	var row string
	err := repository.pool.QueryRow(context.Background(), addApplicationQuery,
		application.Ticker, application.Cost).Scan(&row)
	if err != nil {
		return err
	}

	return nil
}

const selectApplicationsQuery = `select id, ticker, cost, created_at from applications`

func (repository *PostgresqlApplicationsReposity) GetAll() ([]domain.Application, error) {
	rows, err := repository.pool.Query(context.Background(), selectApplicationsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return repository.scanApplications(rows)
}

const selectApplicationsByTicketQuery = `select id, ticker, cost, created_at from applications where ticker=$1`

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
		err := rows.Scan(&app.Id, &app.Ticker, &app.Cost, &app.CreatedAt)
		if err != nil {
			return nil, err
		}

		applications = append(applications, app)
	}

	return applications, nil
}
