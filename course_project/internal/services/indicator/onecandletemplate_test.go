package indicator

import (
	"testing"
	"time"

	"github.com/willsem/tfs-go-hw/course_project/internal/domain"
)

var timeNow = time.Now()

var testingCases = []domain.TickerInfo{
	{
		ProductId: "test",
		Candle: domain.Candle{
			Time:  domain.CandleTime(timeNow),
			Open:  10,
			Close: 10,
		},
	},

	{
		ProductId: "test",
		Candle: domain.Candle{
			Time:  domain.CandleTime(timeNow),
			Open:  10,
			Close: 10,
		},
	},

	{
		ProductId: "test",
		Candle: domain.Candle{
			Time:  domain.CandleTime(timeNow.Add(1 * time.Second)),
			Open:  10,
			Close: 5,
		},
	},

	{
		ProductId: "test",
		Candle: domain.Candle{
			Time:  domain.CandleTime(timeNow.Add(2 * time.Second)),
			Open:  5,
			Close: 10,
		},
	},

	{
		ProductId: "test",
		Candle: domain.Candle{
			Time:  domain.CandleTime(timeNow.Add(2 * time.Second)),
			Open:  5,
			Close: 10,
		},
	},
}

var expectedCases = []domain.Decision{
	domain.NothingDecision,
	domain.NothingDecision,
	domain.BuyDecision,
	domain.SellDecision,
	domain.NothingDecision,
}

func TestOneCandleTemplateOk(t *testing.T) {
	template := NewOneCandleTemplate()

	for testNumber, testCase := range testingCases {
		decision := template.MakeDecision(testCase)

		if decision != expectedCases[testNumber] {
			t.Errorf("Number of case: %d. Expected: %s, Comes: %s.",
				testNumber+1,
				expectedCases[testNumber],
				decision)
		}
	}
}
