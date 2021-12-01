package indicator

import (
	"testing"
	"time"

	"github.com/willsem/tfs-go-hw/course_project/internal/domain"
)

var testingCasesTriple = []domain.TickerInfo{
	{
		ProductId: "test",
	},

	{
		ProductId: "test",
		Candle: domain.Candle{
			Time:  domain.CandleTime(timeNow),
			Open:  5,
			Close: 10,
		},
	},

	{
		ProductId: "test",
		Candle: domain.Candle{
			Time:  domain.CandleTime(timeNow.Add(1 * time.Second)),
			Open:  10,
			Close: 15,
		},
	},

	{
		ProductId: "test",
		Candle: domain.Candle{
			Time:  domain.CandleTime(timeNow.Add(2 * time.Second)),
			Open:  15,
			Close: 20,
		},
	},

	{
		ProductId: "test",
		Candle: domain.Candle{
			Time:  domain.CandleTime(timeNow.Add(3 * time.Second)),
			Open:  20,
			Close: 15,
		},
	},

	{
		ProductId: "test",
		Candle: domain.Candle{
			Time:  domain.CandleTime(timeNow.Add(4 * time.Second)),
			Open:  15,
			Close: 10,
		},
	},

	{
		ProductId: "test",
		Candle: domain.Candle{
			Time:  domain.CandleTime(timeNow.Add(5 * time.Second)),
			Open:  10,
			Close: 5,
		},
	},

	{
		ProductId: "test",
		Candle: domain.Candle{
			Time: domain.CandleTime(timeNow.Add(5 * time.Second)),
		},
	},
}

var expectedCasesTriple = []domain.Decision{
	domain.NothingDecision,
	domain.NothingDecision,
	domain.NothingDecision,
	domain.SellDecision,
	domain.NothingDecision,
	domain.NothingDecision,
	domain.BuyDecision,
	domain.NothingDecision,
}

func TestTripleCandleTemplateOk(t *testing.T) {
	template := NewTripleCandlesTemplate()

	for testNumber, testCase := range testingCasesTriple {
		decision := template.MakeDecision(testCase)

		if decision != expectedCasesTriple[testNumber] {
			t.Errorf("Number of case: %d. Expected: %s, Comes: %s.",
				testNumber+1,
				expectedCasesTriple[testNumber],
				decision)
		}
	}
}
