package subscribe

type candleType string

type SubscribeService interface {
	GetChan() <-chan TickerInfo
	Subscribe(ticker string, candle candleType) error
	Unsubscribe(ticker string, candle candleType) error
}
