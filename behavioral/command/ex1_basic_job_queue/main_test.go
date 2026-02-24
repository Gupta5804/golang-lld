package main

import "testing"

type MockEmailService struct{
	emails []string
}
func (m *MockEmailService) SendEmail(to,body string) error{
	m.emails = append(m.emails,to)
	return nil
}
func TestJobQueue_ProcessAll(t *testing.T) {
	jq := &JobQueue{}
	mock := &MockEmailService{}
	command1 := &EmailCommand{
		emailservice: mock,
		to:"admin@example.com",
		body:"example",
	}
	command2 := &EmailCommand{
		emailservice: mock,
		to:"user@example.com",
		body:"example",
	}
	jq.AddJob(command1)
	jq.AddJob(command2)
	jq.ProcessAll()

	email1_present,email2_present:=false,false
	for _,email := range mock.emails{
		if email == "admin@example.com"{
			email1_present = true
		}
		if email == "user@example.com"{
			email2_present = true
		}
	}
	if !email1_present{
		t.Errorf("Email %s not recorded","admin@example.com")
	}
	if !email2_present{
		t.Errorf("Email %s not recorded", "user@example.com")
	}
}
