package repositories

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"github.com/willsem/tfs-go-hw/course_project/internal/domain"
)

type mockPool struct {
	conn pgxmock.PgxConnIface
}

func newMockPool(conn pgxmock.PgxConnIface) *mockPool {
	return &mockPool{
		conn: conn,
	}
}

func (pool *mockPool) QueryRow(ctx context.Context, query string, fields ...interface{}) pgx.Row {
	return pool.conn.QueryRow(ctx, query, fields)
}

func (pool *mockPool) Query(ctx context.Context, query string, fields ...interface{}) (pgx.Rows, error) {
	return pool.conn.Query(ctx, query, fields)
}

const (
	addApplicationQueryRegexp     = `insert into applications`
	selectApplicationsQueryRegexp = `select`
)

func TestApplicationsRepositoryAddOk(t *testing.T) {
	t.Parallel()

	data := domain.Application{
		Ticker: "APPL",
		Cost:   150,
		Size:   1,
		Type:   domain.BuyAppType,
	}

	mock, err := pgxmock.NewConn()
	assert.Nil(t, err)
	defer mock.Close(context.Background())

	mock.ExpectQuery(addApplicationQueryRegexp).
		WithArgs([]interface{}{data.Ticker, data.Cost, data.Size, data.Type}).
		WillReturnRows(pgxmock.NewRows([]string{"ticker"}).AddRow(data.Ticker))

	pool := newMockPool(mock)
	repo := NewApplicaitionsRepository(pool)
	err = repo.Add(data)

	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestApplicationsRepositoryAddError(t *testing.T) {
	t.Parallel()

	data := domain.Application{
		Ticker: "APPL",
		Cost:   150,
		Size:   1,
		Type:   domain.BuyAppType,
	}

	mock, err := pgxmock.NewConn()
	assert.Nil(t, err)
	defer mock.Close(context.Background())

	dbError := errors.New("test connection crushed")
	mock.ExpectQuery(addApplicationQueryRegexp).
		WithArgs([]interface{}{data.Ticker, data.Cost, data.Size, data.Type}).
		WillReturnError(dbError)

	pool := newMockPool(mock)
	repo := NewApplicaitionsRepository(pool)
	err = repo.Add(data)

	assert.ErrorIs(t, err, dbError)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestApplicationsRepositoryGetAllOk(t *testing.T) {
	t.Parallel()

	data := []domain.Application{
		{
			Id:        1,
			Ticker:    "APPL",
			Cost:      150,
			Size:      1,
			Type:      domain.BuyAppType,
			CreatedAt: time.Time{},
		},
		{
			Id:        2,
			Ticker:    "APPL",
			Cost:      151,
			Size:      2,
			Type:      domain.SellAppType,
			CreatedAt: time.Time{},
		},
	}

	mock, err := pgxmock.NewConn()
	assert.Nil(t, err)
	defer mock.Close(context.Background())

	rows := pgxmock.NewRows([]string{"id", "ticker", "cost", "size", "created_at", "type"})
	for _, row := range data {
		rows = rows.AddRow(row.Id, row.Ticker, row.Cost, row.Size, row.CreatedAt, row.Type)
	}
	mock.ExpectQuery(selectApplicationsQueryRegexp).
		WillReturnRows(rows)

	pool := newMockPool(mock)
	repo := NewApplicaitionsRepository(pool)
	result, err := repo.GetAll()

	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, data, result)
}

func TestApplicationsRepositoryGetAllError(t *testing.T) {
	t.Parallel()

	mock, err := pgxmock.NewConn()
	assert.Nil(t, err)
	defer mock.Close(context.Background())

	dbError := errors.New("test connection crushed")
	mock.ExpectQuery(selectApplicationsQueryRegexp).
		WillReturnError(dbError)

	pool := newMockPool(mock)
	repo := NewApplicaitionsRepository(pool)
	_, err = repo.GetAll()

	assert.Nil(t, mock.ExpectationsWereMet())
	assert.ErrorIs(t, err, dbError)
}

func TestApplicationsRepositoryGetByTickerOk(t *testing.T) {
	t.Parallel()

	data := []domain.Application{
		{
			Id:        1,
			Ticker:    "APPL",
			Cost:      150,
			Size:      1,
			Type:      domain.BuyAppType,
			CreatedAt: time.Time{},
		},
		{
			Id:        2,
			Ticker:    "APPL",
			Cost:      151,
			Size:      2,
			Type:      domain.SellAppType,
			CreatedAt: time.Time{},
		},
	}

	mock, err := pgxmock.NewConn()
	assert.Nil(t, err)
	defer mock.Close(context.Background())

	rows := pgxmock.NewRows([]string{"id", "ticker", "cost", "size", "created_at", "type"})
	for _, row := range data {
		rows = rows.AddRow(row.Id, row.Ticker, row.Cost, row.Size, row.CreatedAt, row.Type)
	}
	mock.ExpectQuery(selectApplicationsQueryRegexp).
		WithArgs([]interface{}{"APPL"}).
		WillReturnRows(rows)

	pool := newMockPool(mock)
	repo := NewApplicaitionsRepository(pool)
	result, err := repo.GetByTicker("APPL")

	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, data, result)
}

func TestApplicationsRepositoryGetByTickerError(t *testing.T) {
	t.Parallel()

	mock, err := pgxmock.NewConn()
	assert.Nil(t, err)
	defer mock.Close(context.Background())

	dbError := errors.New("test connection crushed")
	mock.ExpectQuery(selectApplicationsQueryRegexp).
		WithArgs([]interface{}{"APPL"}).
		WillReturnError(dbError)

	pool := newMockPool(mock)
	repo := NewApplicaitionsRepository(pool)
	_, err = repo.GetByTicker("APPL")

	assert.Nil(t, mock.ExpectationsWereMet())
	assert.ErrorIs(t, err, dbError)
}
