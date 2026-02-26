package main

import (
	"errors"
	"slices"
)

type Request map[string]string

type Handler func(r Request) error

type Middleware func(next Handler) Handler

func RequireAuth(next Handler) Handler {
	return func(r Request) error {
		if r["Authorization"] != "Bearer 123" {
			return errors.New("unauthorized")
		}
		return next(r)
	}
}
func RequireJSON(next Handler) Handler {
	return func(r Request) error {
		if r["Content-Type"] != "application/json" {
			return errors.New("invalid content type")
		}
		return next(r)
	}
}

func BuildChain(final Handler, middlewares ...Middleware) Handler {
	current := final
	slices.Reverse(middlewares)
	for _, middleware := range middlewares {
		current = middleware(current)
	}
	return current
}
