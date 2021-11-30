package repositories

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type IPgxPool interface {
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
}
