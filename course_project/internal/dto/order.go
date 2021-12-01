package dto

import (
	"fmt"
	"net/url"
)

type Order struct {
	OrderType  string
	Symbol     string
	Side       string
	Size       uint64
	LimitPrice float64
}

func (order Order) GetPostData() string {
	v := url.Values{}

	if order.LimitPrice != 0 {
		v.Add("limitPrice", fmt.Sprintf("%f", order.LimitPrice))
	}

	v.Add("orderType", order.OrderType)
	v.Add("symbol", order.Symbol)
	v.Add("side", order.Side)
	v.Add("size", fmt.Sprintf("%d", order.Size))

	return v.Encode()
}
