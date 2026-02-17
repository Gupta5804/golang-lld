package main

import (
	"testing"
)

type MockMessageSender struct {
	sentSubject string
	sentBody    string
}

func (m *MockMessageSender) SendMessage(subject, body string) {
	m.sentSubject = subject
	m.sentBody = body
}
func TestUrgentNotification_AddsPrefix(t *testing.T) {

	mock := &MockMessageSender{}

	urgent := &UrgentNotification{sender: mock}

	urgent.Notify("Warning", "Server Down")

	if mock.sentSubject != "[URGENT] Warning" {
		t.Errorf("Expected subject '[URGENT] Warning', instead got %s", mock.sentSubject)
	}
	if mock.sentBody != "Server Down" {
		t.Errorf("Expected body 'Server Down', got '%s'", mock.sentBody)
	}

}
