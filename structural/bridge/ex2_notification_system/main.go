package main

import "fmt"

type MessageSender interface {
	SendMessage(subject, body string)
}

type EmailSender struct{}

func (e *EmailSender) SendMessage(subject, body string) {
	fmt.Printf("Sending Email:[%s] %s", subject, body)
}

type SMSSender struct{}

func (s *SMSSender) SendMessage(subject, body string) {
	fmt.Printf("Sending SMS:[%s] %s", subject, body)
}

type Notification interface {
	Notify(message, body string)
}

type NormalNotification struct{
	sender MessageSender
}
func (n *NormalNotification) Notify(message, body string) {
	n.sender.SendMessage(message,body)
}

type UrgentNotification struct{
	sender MessageSender
}
func (u *UrgentNotification) Notify(message,body string) {
	subject := fmt.Sprintf("%s %s","[URGENT]",message)
	u.sender.SendMessage(subject,body)
}