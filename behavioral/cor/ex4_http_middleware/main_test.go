package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireAPIKey(t *testing.T) {
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})
	chain := RequireAPIKey(dummyHandler)

	req1 := httptest.NewRequest("GET","/",nil)
	req1.Header.Set("X-API-Key","secret123")
	rr1 := httptest.NewRecorder()
	chain.ServeHTTP(rr1,req1)
	if rr1.Code != http.StatusOK {
		t.Errorf("Expected 200 ok, got %v",rr1.Code)
	}

	req2 := httptest.NewRequest("GET","/",nil)
	req2.Header.Set("X-API-Key","secret321")
	rr2 := httptest.NewRecorder()
	chain.ServeHTTP(rr2,req2)
	if rr2.Code != http.StatusUnauthorized{
		t.Errorf("Expected error 401, got %v",rr2.Code)
	}
}
