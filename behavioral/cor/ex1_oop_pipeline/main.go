package main

import "errors"

type Request struct {
	Headers     map[string]string
	Body        string
	IsProcessed bool
}

// handler interface (so that we can chain)
type Handler interface {
	SetNext(h Handler) Handler
	Execute(r *Request) error
}

// base handler to avoid too much Boilerplate
type BaseHandler struct {
	next Handler
}

func (h *BaseHandler) SetNext(next Handler) Handler {
	h.next = next
	return next
}
func (h *BaseHandler) executeNext(r *Request) error {
	if h.next != nil {
		return h.next.Execute(r)
	}
	return nil // end of chain
}

// concrete handlers
type AuthHandler struct {
	BaseHandler
}

func (a *AuthHandler) Execute(r *Request) error {
	if r.Headers["Authorization"] != "Bearer 123" {
		return errors.New("unauthorized")
	}
	return a.executeNext(r)
}

type JSONHandler struct {
	BaseHandler
}

func (j *JSONHandler) Execute(r *Request) error {
	if r.Headers["Content-Type"] != "application/json" {
		return errors.New("invalid content type")
	}
	return j.executeNext(r)
}

type FinalHandler struct {
	BaseHandler
}

func (f *FinalHandler) Execute(r *Request) error {
	r.IsProcessed = true
	return f.executeNext(r)
}
