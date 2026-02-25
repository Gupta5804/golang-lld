package main

import (
	"context"
	"time"
)

type AsyncCommand interface {
	Execute(ctx context.Context) error
}

type SleepCommand struct {
	Duration time.Duration
}

func (s *SleepCommand) Execute(ctx context.Context) error{
	select{
	case <- time.After(s.Duration):
		return nil
	case <- ctx.Done():
		return ctx.Err()
	}
	//return nil
}

type CommandRunner struct{
	commands []AsyncCommand
}
func (cr *CommandRunner) Add(a AsyncCommand){
	cr.commands = append(cr.commands,a)
}
func (cr *CommandRunner) RunWithTimeout(t time.Duration) error{
	ctx,cancel := context.WithTimeout(context.Background(),t)
	defer cancel()
	for _,command := range cr.commands{
		if err := command.Execute(ctx); err != nil{
			return err
		}
	}
	return nil
}