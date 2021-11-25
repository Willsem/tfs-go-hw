package tradingdto

import (
	"fmt"
)

type Order struct {
	OrderType  string
	Symbol     string
	Side       string
	Size       uint64
	LimitPrice float64
}

func (order Order) GetPostData() string {
	return "orderType=" + order.OrderType + "&" +
		"symbol=" + order.Symbol + "&" +
		"side=" + order.Side + "&" +
		"size=" + fmt.Sprintf("%d", order.Size) + "&" +
		"limitPrice=" + fmt.Sprintf("%f", order.LimitPrice)
}
