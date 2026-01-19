package main

import (
	"fmt"
	"testing"
)

type MockNotifier struct {
	LastMessage string
	LastUserID int
	CallCount int
}

func (m *MockNotifier) Send(message string, userID int) {
	m.LastMessage = message
	m.LastUserID = userID
	m.CallCount++
}

func TestNotificationService_NotifyUser(t *testing.T) {
	// setup
	mock := &MockNotifier{}
	service := NewNotificationService(mock)

	// run the service 
	service.NotifyUser(22,"test")
	// assert
	if mock.CallCount != 1 {
		t.Errorf("Expected 1 Call, got %d", mock.CallCount)
	}
	if mock.LastUserID != 22 {
		t.Errorf("Expected userID 22, got %d",mock.LastUserID)
	}
	if mock.LastMessage != "test" {
		t.Errorf("Expected message 'test', got %s", mock.LastMessage)
	}

	fmt.Println("Test Passed !")
}