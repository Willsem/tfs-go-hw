package applications

import "github.com/willsem/tfs-go-hw/course_project/internal/domain"

type ApplicationsRepository interface {
	Add(application domain.Application) error
	GetAll() ([]domain.Application, error)
	GetByTicket(ticker string) ([]domain.Application, error)
}
