package feed_test

import (
	"testing"
	"time"

	"github.com/Gupta5804/golang-lld/systems/phase1/03_stock_alerts/domain"
	"github.com/Gupta5804/golang-lld/systems/phase1/03_stock_alerts/feed"
)

func TestEngine_UpdatePrice_Evaluations(t *testing.T) {
	tests := []struct {
		name               string
		setup              func() (*feed.Engine, <-chan domain.Command)
		triggerSymbol      string
		triggerPrices      []int64
		expectedCommandCount int
	}{
		{
			name: "The `Above` trigger",
			setup: func() (*feed.Engine, <-chan domain.Command) {
				commandChan := make(chan domain.Command, 1)
				engine := feed.NewEngine(commandChan)
				alertCondition := domain.AlertCondition{
					AlertID:      "alert1",
					UserID:       "user1",
					TickerSymbol: "AAPL",
					TargetPrice:  10000,
					Direction:    domain.Above,
				}
				engine.Subscribe(alertCondition)
				return engine, commandChan
			},
			triggerSymbol:      "AAPL",
			triggerPrices:      []int64{10500},
			expectedCommandCount: 1,
		},
		{
			name: "The `Below` trigger",
			setup: func() (*feed.Engine, <-chan domain.Command) {
				commandChan := make(chan domain.Command, 1)
				engine := feed.NewEngine(commandChan)
				alertCondition := domain.AlertCondition{
					AlertID:      "alert1",
					UserID:       "user1",
					TickerSymbol: "AAPL",
					TargetPrice:  10000,
					Direction:    domain.Below,
				}
				engine.Subscribe(alertCondition)
				return engine, commandChan
			},
			triggerSymbol:      "AAPL",
			triggerPrices:      []int64{9500},
			expectedCommandCount: 1,
		},
		{
			name: "The `No Trigger` trigger",
			setup: func() (*feed.Engine, <-chan domain.Command) {
				commandChan := make(chan domain.Command, 1)
				engine := feed.NewEngine(commandChan)
				alertCondition := domain.AlertCondition{
					AlertID:      "alert1",
					UserID:       "user1",
					TickerSymbol: "AAPL",
					TargetPrice:  10000,
					Direction:    domain.Above,
				}
				engine.Subscribe(alertCondition)
				return engine, commandChan
			},
			triggerSymbol:      "AAPL",
			triggerPrices:      []int64{9500},
			expectedCommandCount: 0,
		},
		{
			name: "The `One-Shot` Constraint",
			setup: func() (*feed.Engine, <-chan domain.Command) {
				commandChan := make(chan domain.Command, 2)
				engine := feed.NewEngine(commandChan)
				alertCondition := domain.AlertCondition{
					AlertID:      "alert1",
					UserID:       "user1",
					TickerSymbol: "AAPL",
					TargetPrice:  10000,
					Direction:    domain.Above,
				}
				engine.Subscribe(alertCondition)
				return engine, commandChan
			},
			triggerSymbol:      "AAPL",
			triggerPrices:      []int64{10500, 10600},
			expectedCommandCount: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			engine, ch := tc.setup()

			for _, price := range tc.triggerPrices {
				engine.UpdatePrice(tc.triggerSymbol, price)
			}

			var actualCommands int
			done := false

			for !done {
				select {
				case <-ch:
					actualCommands++
				case <-time.After(10 * time.Millisecond):
					done = true
				}
			}

			if actualCommands != tc.expectedCommandCount {
				t.Errorf("Expected %d commands, got %d",tc.expectedCommandCount,actualCommands)
			}
		})
	}
}
