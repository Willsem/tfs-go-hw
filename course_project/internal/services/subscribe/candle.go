package subscribe

import (
	"strconv"
	"strings"
	"time"
)

type Candle struct {
	Time   candleTime  `json:"time"`
	Open   candleFloat `json:"open"`
	High   candleFloat `json:"high"`
	Low    candleFloat `json:"low"`
	Close  candleFloat `json:"close"`
	Volume int64       `json:"volume"`
}

type candleTime time.Time

func (t *candleTime) UnmarshalJSON(data []byte) error {
	millis, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = candleTime(time.Unix(0, millis*int64(time.Millisecond)))
	return nil
}

type candleFloat float64

func (f *candleFloat) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}
	*f = candleFloat(value)
	return nil
}
