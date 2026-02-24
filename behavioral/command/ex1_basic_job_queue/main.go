package main

// command interface
type Command interface {
	Execute() error
}

// chef / receiver
// interface
type EmailService interface {
	SendEmail(to, body string) error
}

// no need for struct of this interface as we are just mocking it in test

// command struct/ticket
type EmailCommand struct {
	emailservice EmailService
	to           string
	body         string
}
func (ec *EmailCommand) Execute() error{
	return ec.emailservice.SendEmail(ec.to,ec.body)
}
// invoker/waiter
type JobQueue struct{
	commands []Command
}
func(jq *JobQueue) AddJob(command Command){
	jq.commands = append(jq.commands,command)
}
func(jq *JobQueue) ProcessAll() error {
	for _,command := range jq.commands{
		if err:=command.Execute();err!=nil{
			return err
		}
	}
	jq.commands = nil
	return nil
}
