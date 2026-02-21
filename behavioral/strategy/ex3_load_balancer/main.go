package main

import (
	"errors"
	"math/rand"
)

type RoutingStrategy interface {
	Route(backends []string) string
}

type RoundRobinStrategy struct {
	currIndex int
}

func (r *RoundRobinStrategy) Route(backends []string) string {
	curr := backends[r.currIndex]
	r.currIndex++
	r.currIndex = r.currIndex % len(backends)
	return curr
}

type RandomStrategy struct{}

func (r *RandomStrategy) Route(backends []string) string {
	randomIndex := rand.Intn(len(backends))
	return backends[randomIndex]
}

type LoadBalancer struct {
	strategy RoutingStrategy
	backends []string
}

func NewLoadBalancer(backends []string, strategy RoutingStrategy) *LoadBalancer {
	return &LoadBalancer{
		strategy: strategy,
		backends: backends,
	}
}
func (lb *LoadBalancer) GetNextServer() (string, error) {
	if len(lb.backends) == 0 {
		return "", errors.New("no servers")
	}
	return lb.strategy.Route(lb.backends),nil
}
