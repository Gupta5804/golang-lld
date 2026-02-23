package main

import (
	"slices"
	"sync"
)

type Subscriber chan int

type MetricsBroker struct {
	subscribers map[string][]Subscriber
	mu          sync.RWMutex
}

func NewMetricsBroker() *MetricsBroker {
	return &MetricsBroker{
		subscribers: make(map[string][]Subscriber),
	}
}

func (mb *MetricsBroker) Subscribe(metric string) Subscriber {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	ch := make(Subscriber, 1)
	mb.subscribers[metric] = append(mb.subscribers[metric], ch)
	return ch
}
func (mb *MetricsBroker) Unsubscribe(metric string, subscriber Subscriber) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	for i, ch := range mb.subscribers[metric] {
		if subscriber == ch {
			mb.subscribers[metric] = slices.Delete(mb.subscribers[metric],i,i+1)
			close(ch)
			break
		}
	}
}
func (mb *MetricsBroker) Publish(metric string,val int) {
	mb.mu.RLock()
	defer mb.mu.RUnlock()
	if chans,found := mb.subscribers[metric]; found{
		for _,ch := range chans{
			select{
			case ch<-val:
			default:
			}
		}
	}
}
