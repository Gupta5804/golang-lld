package main

import (
	"fmt"
	"strings"
)

type Job func(input string) string

func BaseJob(input string) string {
	return "Processed: " + input
}

func Auth(j Job, user string) Job {
	return func(input string) string{
		if user == "admin" {
			return j(input)
		} else {
			return "Unauthorized"
		}
	}
}
func Logging(j Job) Job {
	return func(input string) string {
		fmt.Printf("Log: %s\n", input)
		return j(input)
	}
}

func Sanitize(j Job) Job {
	return func(input string) string {
		cleanInput := strings.TrimSpace(input)
		return j(cleanInput)
	}
}
