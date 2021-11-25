package tradingbot

type TradingBot interface {
	Start() error
	Configure(command Command) error
}
