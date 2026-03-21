package domain

import "time"

type Direction int

const (
	Above Direction = iota
	Below
)

func (d Direction) String() string {
	switch d {
	case Above:
		return "ABOVE"
	case Below:
		return "BELOW"
	default:
		return "UNKNOWN"
	}
}

type User struct {
	ID string
}

type Ticker struct {
	Symbol       string
	CurrentPrice int64
}

type AlertCondition struct {
	AlertID      string
	UserID       string
	TickerSymbol string
	TargetPrice  int64
	Direction    Direction
}

type AlertSnapshot struct {
	Condition      AlertCondition
	TriggeredPrice int64
	Timestamp      time.Time
}

type Command interface {
	Execute() error
}
