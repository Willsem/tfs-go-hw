package telegram

import "github.com/willsem/tfs-go-hw/course_project/internal/domain"

type ApplicationsRepository interface {
	GetAll() ([]domain.Application, error)
	GetByTicker(ticker string) ([]domain.Application, error)
}
