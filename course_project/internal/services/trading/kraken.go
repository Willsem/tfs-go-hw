package trading

import (
	"net/http"

	"github.com/willsem/tfs-go-hw/course_project/internal/config"
)

type KrakenTradingService struct {
	config config.Kraken
	client http.Client
}

func New(config config.Kraken) (*KrakenTradingService, error) {
	return &KrakenTradingService{
		config: config,
		client: http.Client{},
	}, nil
}

func (service *KrakenTradingService) Buy(ticker Ticker) error {
	return nil
}

func (service *KrakenTradingService) Sell(ticker Ticker) error {
	return nil
}
