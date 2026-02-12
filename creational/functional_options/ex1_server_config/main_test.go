package main

import (
	"testing"
)

func TestNewServer_Defaults(t *testing.T) {

	server := NewServer()

	if server.host != "localhost" {
		t.Errorf("Expected 'localhost' as Host instead got %s", server.host)
	}
	if server.port != 8080 {
		t.Errorf("Expected port as 8080 instead got %d", server.port)
	}
	if server.timeout != 30 {
		t.Errorf("Expected timeout to be 30s instead got %ds", server.timeout)
	}
	if server.maxconn != 100 {
		t.Errorf("Expected maxconn as 100 , instead got %d", server.maxconn)
	}

}

func TestNewServer_WithOptions(t *testing.T) {
	server := NewServer(
		WithPort(9090),
		WithTimeout(60),
	)

	if server.host != "localhost" {
		t.Errorf("Expected 'localhost' as Host instead got %s", server.host)
	}
	if server.port != 9090 {
		t.Errorf("Expected port as 9090 instead got %d", server.port)
	}
	if server.timeout != 60 {
		t.Errorf("Expected timeout to be 60s instead got %ds", server.timeout)
	}
	if server.maxconn != 100 {
		t.Errorf("Expected maxconn as 100 , instead got %d", server.maxconn)
	}

}
