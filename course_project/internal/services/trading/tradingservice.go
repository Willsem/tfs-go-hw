package trading

type TradingService interface {
	Buy(ticker Ticker) error
	Sell(ticker Ticker) error
}
