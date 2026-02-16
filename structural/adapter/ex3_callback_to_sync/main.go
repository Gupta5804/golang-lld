package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// --- Target Interface (synchronous) ---
type AuthService interface {
	Login(username, password string) error
}

// --- Adaptee: Callback-based 3rd Party Lib ---
type OAuthClient struct{}

// AsyncLogin returns immediately. It simulates network delay,
// then calls EITHER onSuccess OR onFailure
func (o *OAuthClient) AsyncLogin(user, pass string, onSuccess func(), onFailure func(error)) {
	go func() {
		// simulate network latency
		time.Sleep(300 * time.Millisecond)
		if user == "admin" && pass == "secret" {
			onSuccess()
		} else {
			onFailure(errors.New("invalid credentials"))
		}
	}()
}

type LoginAdapter struct {
	oauth   *OAuthClient
	Timeout time.Duration
}

func (l *LoginAdapter) Login(username, password string) error {
	resultChan := make(chan error)

	onSuccess := func() {
		resultChan <- nil
	}

	onFailure := func(err error) {
		resultChan <- err
	}

	ctx, cancel := context.WithTimeout(context.Background(), l.Timeout)
	defer cancel()
	l.oauth.AsyncLogin(username, password, onSuccess, onFailure)

	select {
	case err := <-resultChan:
		return err
	case <-ctx.Done():
		return fmt.Errorf("request timed out")
	}
}
