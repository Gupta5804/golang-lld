package main

import "testing"

type MockAccount struct {
	called bool
}
func (m *MockAccount) CheckKYC(user string) error {
	m.called = true
	return nil
}

type MockBank struct {
	called bool
}
func (m *MockBank) Debit(amount int) error {
	m.called = true
	return nil
}

type MockBlockchain struct {
	called bool
}
func (m *MockBlockchain) Credit(amount int) error {
	m.called = true
	return nil
}

type MockNotification struct {
	called bool
}
func (m *MockNotification) SendReceipt(user string, amount int) {
	m.called = true
}

func TestWalletFacade_BuyCrypto(t *testing.T) {
	account := &MockAccount{}
	bank := &MockBank{}
	blockchain := &MockBlockchain{}
	notification := &MockNotification{}

	facade := &WalletFacade{
		account:      account,
		bank:         bank,
		blockchain:   blockchain,
		notification: notification,
	}

	
	err := facade.BuyCrypto(100, "user1")

	
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}

	if !account.called {
		t.Error("Step 1 Failed: Account CheckKYC was not called")
	}
	if !bank.called {
		t.Error("Step 2 Failed: Bank Debit was not called")
	}
	if !blockchain.called {
		t.Error("Step 3 Failed: Blockchain Credit was not called")
	}
	if !notification.called {
		t.Error("Step 4 Failed: Notification SendReceipt was not called")
	}
}
