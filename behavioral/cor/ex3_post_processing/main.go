package main

import (
	"slices"
	"time"
)

type Request map[string]string
type Handler func(r Request) error
type Middleware func(h Handler) Handler

func BuildChain(final Handler, middlewares ...Middleware) Handler {
	current := final
	slices.Reverse(middlewares)
	for _, middleware := range middlewares {
		current = middleware(current)
	}
	return current
}

func MetricsMiddleware(metrics map[string]string) Middleware {
	return func(h Handler) Handler {
		return func(r Request) error {
			start := time.Now()

			err := h(r)
			duration:= time.Since(start)

			if err != nil {
				metrics["status"] = "fail"
			} else {
				metrics["status"] = "success"
			}

			metrics["duration"] = duration.String()
			return err
		}
	}
}
