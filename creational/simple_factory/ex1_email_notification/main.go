package main

import "fmt"

// -- Concrete Implementation ---
type EmailNotification struct{
	sender string
}

func (e *EmailNotification) Send(message string ,userID int) {
	fmt.Printf("Sending Email to User %d from %s: %s \n",userID, e.sender, message)
}

// --- business logic ---
type NotificationService struct {
	notifier *EmailNotification
}

func (n *NotificationService) NotifyUser(userID int, msg string) {
	n.notifier.Send(msg, userID)
}

func main() {
	emailNotif := &EmailNotification{sender:"admin@app.com"}

	service := &NotificationService{
		notifier : emailNotif,
	}
	service.NotifyUser(42,"Welcome to our platform!")
}