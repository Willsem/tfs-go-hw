package domain

type CandleType string

const (
	Candle1m CandleType = "candles_trade_1m"
	Candle5m CandleType = "candles_trade_5m"
	Candle1h CandleType = "candles_trade_1h"
)
