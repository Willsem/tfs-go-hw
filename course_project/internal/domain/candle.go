package domain

import (
	"strconv"
	"strings"
	"time"
)

type Candle struct {
	Time   CandleTime  `json:"time"`
	Open   CandleFloat `json:"open"`
	High   CandleFloat `json:"high"`
	Low    CandleFloat `json:"low"`
	Close  CandleFloat `json:"close"`
	Volume int64       `json:"volume"`
}

type CandleTime time.Time

func (t *CandleTime) UnmarshalJSON(data []byte) error {
	millis, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = CandleTime(time.Unix(0, millis*int64(time.Millisecond)))
	return nil
}

type CandleFloat float64

func (f *CandleFloat) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}
	*f = CandleFloat(value)
	return nil
}
