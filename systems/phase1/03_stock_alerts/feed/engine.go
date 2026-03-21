package feed

import (
	"fmt"
	"sync"
	"time"

	"github.com/Gupta5804/golang-lld/systems/phase1/03_stock_alerts/domain"
)

type Engine struct {
	subscriptions map[string]map[string]domain.AlertCondition // outer map key : Ticker Symbol ; Inner map key: AlertID
	mu            sync.RWMutex
	commandChan   chan<- domain.Command
}

func NewEngine(commandChan chan<- domain.Command) *Engine {
	return &Engine{
		subscriptions: make(map[string]map[string]domain.AlertCondition),
		commandChan:   commandChan,
	}
}
func (e *Engine) Subscribe(condition domain.AlertCondition) {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, exists := e.subscriptions[condition.TickerSymbol]
	if !exists {
		e.subscriptions[condition.TickerSymbol] = make(map[string]domain.AlertCondition)
	}
	e.subscriptions[condition.TickerSymbol][condition.AlertID] = condition
}

func (e *Engine) Unsubscribe(symbol, alertID string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if innerMap, exists := e.subscriptions[symbol]; exists {
		delete(innerMap, alertID)
	}
}
func (e *Engine) UpdatePrice(symbol string, newPrice int64) {
	e.mu.Lock()
	defer e.mu.Unlock()

	conditions, exists := e.subscriptions[symbol]
	if !exists {
		return
	}
	for alertID, condition := range conditions {
		isTriggered := false
		switch condition.Direction {
		case domain.Above:
			if newPrice >= condition.TargetPrice {
				isTriggered = true
			}
		case domain.Below:
			if newPrice <= condition.TargetPrice {
				isTriggered = true
			}
		}
		if isTriggered {
			cmd := &AlertCommand{
				Snapshot: domain.AlertSnapshot{
					Condition:      condition,
					TriggeredPrice: newPrice,
					Timestamp:      time.Now(),
				},
			}
			select {
			case e.commandChan <- cmd:
			default:
				// channel buffer is full, workers are too slow
				// we drop the alert to save the engine
			}

			delete(conditions, alertID)
		}
	}
}

type AlertCommand struct {
	Snapshot domain.AlertSnapshot
}

func (a *AlertCommand) Execute() error {
	fmt.Printf("[Alert Sent] Ticker: %s | Triggered at: $%.2f | Time: %s/n",
		a.Snapshot.Condition.TickerSymbol,
		float64(a.Snapshot.TriggeredPrice)/100.0,
		a.Snapshot.Timestamp.Format("15:04:05.000"))
	return nil
}
