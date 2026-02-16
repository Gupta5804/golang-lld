package main

import (
	"testing"
	"time"
)

func TestCryptoAdapter_Success(t *testing.T) {
	// setup the adaptee
	processor := &CryptoProcessor{}

	var gateway PaymentGateway = &CryptoAdapter{
		processor: processor,
		Timeout:   2 * time.Second,
	}

	if err := gateway.Pay(100.0); err != nil {
		t.Errorf("Expected success, got this error:%v", err)
	}

}

func TestCryptoAdapter_Timeout(t *testing.T) {
	processor := &CryptoProcessor{}

	var gateway PaymentGateway = &CryptoAdapter{
		processor: processor,
		Timeout:   1 * time.Millisecond,
	}

	err := gateway.Pay(100.0)
	if err == nil {
		t.Error("Expected a timeout error, but got nil(success)")
	}

	if err != nil && err.Error() != "payment timed out" {
		t.Errorf("Expected 'payment timed out' error , got : %v", err)
	}
}
