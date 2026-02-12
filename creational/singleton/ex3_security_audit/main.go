package main

import (
	"sync"
)

type AuditLogger struct {
	count map[string]int
	rwmu  sync.RWMutex
}

func (a *AuditLogger) LogEvent(name, log string) {
	a.rwmu.Lock()
	defer a.rwmu.Unlock()
	a.count[name]++
}

func (a *AuditLogger) GetEventCount(name string) int {
	a.rwmu.RLock()
	defer a.rwmu.RUnlock()

	return a.count[name]
}

var (
	instance *AuditLogger
	once     sync.Once
)

func GetInstance() *AuditLogger {
	once.Do(func() {
		instance = &AuditLogger{
			count: make(map[string]int),
		}
	})
	return instance
}
