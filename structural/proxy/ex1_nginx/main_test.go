package main

import (
	"testing"
)

type ApplicationServer struct {
	callCount int
}

func (a *ApplicationServer) HandleRequest(url, method string) (int, string) {
	a.callCount++
	return 200, "OK"
}
func TestProxy_RateLimiting(t *testing.T) {
	realServer := &ApplicationServer{}
	proxy := &NginxProxy{
		server:       realServer,
		rateLimit:    2,
		requestCount: map[string]int{},
		cache:        map[string]string{},
	}

	code, _ := proxy.HandleRequest("/app", "GET")
	if code != 200 {
		t.Errorf("Request 1 failed. Got %d", code)
	}

	code, _ = proxy.HandleRequest("/app", "GET")
	if code != 200 {
		t.Errorf("Request 2 failed. Got %d", code)
	}

	code, msg := proxy.HandleRequest("/app", "GET")
	if code != 429 {
		t.Errorf("Expected 429 too many requet, got %d", code)
	}
	if msg != "Not Allowed" {
		t.Errorf("Expected 'Not Allowed', got %s", msg)
	}
}

func TestProxy_Caching(t *testing.T) {
	realServer := &ApplicationServer{}

	proxy := &NginxProxy{
		server:       realServer,
		rateLimit:    100,
		cache:        map[string]string{},
		requestCount: map[string]int{},
	}

	proxy.HandleRequest("/video", "GET") // hits realServer
	proxy.HandleRequest("/video", "GET") // hits cache
	proxy.HandleRequest("/video", "GET") // hits cache

	if realServer.callCount != 1 {
		t.Errorf("Expected real server to be called just once, but was called %d times", realServer.callCount)
	}
}
