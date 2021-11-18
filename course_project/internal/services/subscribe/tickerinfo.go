package subscribe

type TickerInfo struct {
	Feed      string `json:"feed"`
	ProductId string `json:"product_id"`
	Candle    Candle `json:"candle"`
}
