package main

import "fmt"

func ErrorMessage(err error) {
	fmt.Printf("Client Error: [ %s ] \n", err)
}

func LogMessage(msg string) {
	fmt.Printf("Client Log: [ %s ] \n", msg)
}
