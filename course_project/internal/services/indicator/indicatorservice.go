package indicator

import "github.com/willsem/tfs-go-hw/course_project/internal/domain"

type IndicatorService interface {
	MakeDesicion(ticker domain.TickerInfo) Decision
}
