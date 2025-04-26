package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// setupSignalHandling configures the signal handling and returns a done channel
// The signal handler now signals the done channel
func setupSignalHandling() chan bool {
	// Create a channel to receive OS signals
	sigs := make(chan os.Signal, 1)
	// Create a channel to signal when the program is done normally
	done := make(chan bool, 1)

	// Register the 'sigs' channel to receive Interrupt (Ctrl+C)
	// and SIGTERM signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine to handle received signals
	go func() {
		// Wait for a signal on the 'sigs' channel
		<-sigs
		fmt.Println("\nOk. See you next time. Goodbye!")
		// Signal the done channel to indicate termination
		done <- true
	}()

	return done
}
