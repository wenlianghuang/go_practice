package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

// main is the entry point of the application. It sets up logging to a file named "app.log".
// If the log file cannot be opened, the application will print an error message and exit.
// The logger is configured to include standard log flags (date and time).
// The function also includes a deferred function to handle any panics that occur,
// logging the panic message and stack trace before the application exits.
// The application logs the start and end of the execution, as well as the result of someFunction.
// main is the entry point of the application.
// It opens a log file for writing logs, sets up a logger, and handles any panics that occur during execution.
func main() {
	// Open the log file "app.log" with create, write-only, and append modes, and set permissions to 0644.
	logfile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		// Print an error message and exit if the log file cannot be opened.
		fmt.Println("Failed to open log file:", err)
		os.Exit(1)
	}

	// Create a new logger that writes to the log file with standard log flags.
	logger := log.New(logfile, "", log.LstdFlags)

	// Defer a function to handle any panics that occur.
	defer func() {
		if r := recover(); r != nil {
			// Log the panic message and stack trace if a panic occurs.
			logger.Println("Panic occurred:", r)
			logger.Println(string(debug.Stack()))
		}
	}()

	// Log the start of the application.
	logger.Println("Starting application...")

	// Call someFunction and log its result.
	result := someFunction()
	logger.Println("Result:", result)

	// Log the successful completion of the application.
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
