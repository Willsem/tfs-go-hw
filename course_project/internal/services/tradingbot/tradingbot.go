package tradingbot

type TradingBot interface {
	Start() error

	Continue() error
	Pause() error

	AddTicker(ticker string) error
	RemoveTicker(ticker string) error
}
