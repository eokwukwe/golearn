package main

import (
	"fmt"
)

func main() {

	// Set up signal handling
	done := setupSignalHandling()

	// Create and start a new game
	game := NewGame(1, 100)
	fmt.Println("Welcome to the Guess Game!")
	fmt.Println("Guess a number between 1 and 100. Enter <quit> or use ctrl+c to exit.")

	// Start the game loop in a goroutine so main can wait on the 'done' channel
	go game.Start(done)

	// Wait for the game to finish or a signal to be received
	<-done
}
