package main

import "fmt"

func ErrorMessage(err error) {
	fmt.Printf("Server Error: [ %s ] \n", err)
}

func LogMessage(msg string) {
	fmt.Printf("Server Log: [ %s ] \n", msg)
}
