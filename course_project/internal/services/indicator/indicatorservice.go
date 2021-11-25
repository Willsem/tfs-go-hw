package indicator

import "github.com/willsem/tfs-go-hw/course_project/internal/domain"

type IndicatorService interface {
	MakeDecision(ticker domain.TickerInfo) Decision
}
