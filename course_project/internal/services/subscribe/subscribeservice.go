package subscribe

type SubscribeService interface {
	GetChan() <-chan TickerInfo
	Subscribe(ticker string) error
	Unsubscribe(ticker string) error
}
