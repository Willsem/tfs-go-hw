POSTRESQL_CONNECTION_STRING := "postgres://postgres:123456@localhost:5432/trading_db?sslmode=disable"


build:
	go build -v ./cmd/api

run:
	go run -v ./cmd/api

debug:
	go run -v ./cmd/api -debug

test:
	go test ./...

coverage:
	go test -cover ./...

migrateup:
	migrate -path db/migration -database $(POSTRESQL_CONNECTION_STRING) -verbose up

migratedown:
	migrate -path db/migration -database $(POSTRESQL_CONNECTION_STRING) -verbose down


.DEFAULT_GOAL := build
