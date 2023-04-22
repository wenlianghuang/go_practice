package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

func main() {
	logfile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		os.Exit(1)
	}

	logger := log.New(logfile, "", log.LstdFlags)

	defer func() {
		if r := recover(); r != nil {
			logger.Println("Panic occurred:", r)
			logger.Println(string(debug.Stack()))
		}
	}()

	logger.Println("Starting application...")
	result := someFunction()
	logger.Println("Result:", result)
	logger.Println("Application finished successfully.")
}

func someFunction() int {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Sprintf("Unexpected error: %v", r))
		}
	}()

	// Some code that may cause panic
	//panic("Something went wrong!")
	panic(1)
}
