package handlers

type TradingBot interface {
	IsWorking() bool
	Continue() error
	Pause() error

	Tickers() []string
	AddTicker(ticker string) error
	RemoveTicker(ticker string) error

	ChangeSize(newSize uint64)
}
