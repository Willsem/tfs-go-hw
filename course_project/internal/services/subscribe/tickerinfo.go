package subscribe

type TickerInfo struct {
	Time     uint64  `json:"time"`
	Name     string  `json:"product_id"`
	BidPrice float64 `json:"bid"`
	AskPrice float64 `json:"ask"`
	Change   float64 `json:"change"`
}
