package applications

import (
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

func (repository *PostgresqlApplicationsReposity) Add(application domain.Application) error {
	return nil
}

func (repository *PostgresqlApplicationsReposity) GetAll() ([]domain.Application, error) {
	return nil, nil
}

func (repository *PostgresqlApplicationsReposity) GetByTicket(ticker string) ([]domain.Application, error) {
	return nil, nil
}
