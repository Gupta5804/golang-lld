package main

import "fmt"

// --- behaviour/contract ---
type Notifier interface {
	Send(message string, userID int)
}

// -- Concrete Implementation EMAIL---
type EmailNotification struct {
	sender string
}

func (e *EmailNotification) Send(message string, userID int) {
	fmt.Printf("Sending Email to User %d from %s: %s \n", userID, e.sender, message)
}

// -- Concrete Implementation SMS ---
type SMSNotification struct {
	sender string
}

func (s *SMSNotification) Send(message string, userID int) {
	fmt.Printf("Sending SMS to User %d from %s: %s \n", userID, s.sender, message)
}

// --- business logic ---
type NotificationService struct {
	notifier Notifier
}

func (n *NotificationService) NotifyUser(userID int, msg string) {
	n.notifier.Send(msg, userID)
}

// Constructor function / builder function
func NewNotificationService(n Notifier) *NotificationService {
	return &NotificationService{
		notifier: n,
	}
}

func main() {
	emailNotif := &EmailNotification{sender: "admin@app.com"}
	smsNotif := &SMSNotification{sender: "admin2@app.com"}
	service1 := NewNotificationService(emailNotif)
	service2 := NewNotificationService(smsNotif)

	service1.NotifyUser(42, "Welcome to our platform!")
	service2.NotifyUser(24, "hello from the other side")
}
