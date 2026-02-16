package main

import (
	"testing"
	"time"
)

func TestLoginAdapter_Success(t *testing.T) {
	oauth := &OAuthClient{}

	user, pass := "admin", "secret"

	var loginService AuthService = &LoginAdapter{
		oauth:   oauth,
		Timeout: 1 * time.Second,
	}

	if err := loginService.Login(user, pass); err != nil {
		t.Errorf("Expected success, got %v", err)
	}

}

func TestLoginAdapter_Timeout(t *testing.T) {

	oauth := &OAuthClient{}

	user, pass := "admin", "secret"

	var loginService AuthService = &LoginAdapter{
		oauth:   oauth,
		Timeout: 100 * time.Millisecond,
	}

	err := loginService.Login(user, pass)

	if err == nil {
		t.Error("Expected error here , but got success")
	}
	if err != nil && err.Error() != "request timed out" {
		t.Errorf("Expected 'request timed out error', but got %v", err)
	}
}
