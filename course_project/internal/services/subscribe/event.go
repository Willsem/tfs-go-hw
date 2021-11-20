package subscribe

type event struct {
	Event      string   `json:"event,omitempty"`
	Feed       string   `json:"feed,omitempty"`
	ProductIds []string `json:"product_ids,omitempty"`
	ProductId  string   `json:"product_id,omitempty"`
	Message    string   `json:"message,omitempty"`
	Candle     Candle   `json:"candle,omitempty"`
}
