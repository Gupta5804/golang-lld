package main

import (
	"errors"
	"testing"
)

func TestCache_Hit(t *testing.T) {
	dbCalled := false

	realDB := func(key string) (string, error) {
		dbCalled = true
		return "real_db_data", nil
	}

	cache := map[string]string{
		"user:1": "cached_data",
	}

	cachedFetcher := WithCache(realDB, cache)

	val, err := cachedFetcher("user:1")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if val != "cached_data" {
		t.Errorf("Expected 'cached_data', got %s", val)
	}
	if dbCalled {
		t.Error("Cache hit should not call the underlying DB, but it did")
	}
}
func TestRetry_Success(t *testing.T) {
	attempts := 0
	flakyDB := func(key string) (string, error) {
		attempts++
		if attempts < 3 {
			return "", errors.New("connection reset")
		}
		return "success_data", nil
	}
	retryFetcher := WithRetry(flakyDB, 3)

	val, err := retryFetcher("any_key")

	if err != nil {
		t.Fatalf("Expected success after retries, got error: %v", err)
	}
	if val != "success_data" {
		t.Errorf("Expected 'success_data, got '%s'", val)
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestRetry_Fail(t *testing.T) {
	brokenDB := func(key string) (string, error) {
		return "", errors.New("persistent failure")
	}

	retryFetcher := WithRetry(brokenDB, 2)

	_, err := retryFetcher("key1")

	if err == nil {
		t.Error("Expected error after max retries, got success instead")
	}
}
