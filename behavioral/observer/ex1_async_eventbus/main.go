package main

import "sync"

type Subscriber chan string

type EventBus struct {
	subscribers map[string][]Subscriber
	mu          sync.RWMutex
}

func NewEventBus() *EventBus{
	return &EventBus{
		subscribers:make(map[string][]Subscriber),
	}
}

func (eb *EventBus) Subscribe(eventType string) Subscriber{
	eb.mu.Lock()
	defer eb.mu.Unlock()

	ch := make(Subscriber)
	eb.subscribers[eventType] = append(eb.subscribers[eventType],ch)
	return ch
}

func (eb *EventBus) Publish(eventType string,msg string) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if chans,found := eb.subscribers[eventType];found{
		for _,ch := range chans{
			go func(subscriber Subscriber) {
				subscriber <- msg
			}(ch)
		}
	} 
}
