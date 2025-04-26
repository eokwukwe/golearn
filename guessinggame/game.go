package main

import (
	"fmt"
	"math/rand"
)

// Game represents the guessing game state
type Game struct {
	minRange     int
	maxRange     int
	secretNumber int
}

// NewGame creates a new game instance with the given range
func NewGame(min, max int) *Game {
	return &Game{
		minRange:     min,
		maxRange:     max,
		secretNumber: rand.Intn(max-min+1) + min,
	}
}

// Start begins the game loop
func (g *Game) Start(done chan bool) {
	defer func() {
		// Ensure the done channel is signaled when the function exits
		// This handles both normal game end and potential panics (though unlikely here)
		done <- true
	}()

	inputHandler := NewInputHandler()

	for {
		input := inputHandler.ReadInput()

		if input.isEmpty {
			fmt.Println("You did not enter anything. Enter a number.")
			continue
		}

		if input.isQuit {
			fmt.Println("\nThanks for playing. Goodbye!")
			return // Exit the goroutine
		}

		if input.err != nil {
			fmt.Printf("Invalid input: %v. Please enter a whole number.\n", input.err)
			continue
		}

		if !g.isValidGuess(input.value) {
			fmt.Printf("Your guess must be between %d and %d.\n", g.minRange, g.maxRange)
			continue
		}

		if g.checkGuess(input.value) {
			return
		}
	}
}

// isValidGuess checks if the guess is within the valid range
func (g *Game) isValidGuess(guess int) bool {
	return guess >= g.minRange && guess <= g.maxRange
}

// checkGuess evaluates the guess and returns true if correct
func (g *Game) checkGuess(guess int) bool {
	switch {
	case guess < g.secretNumber:
		fmt.Println("Too low!")
		return false
	case guess > g.secretNumber:
		fmt.Println("Too high!")
		return false
	default: // guess == g.secretNumber
		fmt.Println("Congratulations! 🎉 You guessed the number!")
		fmt.Println("Thank you for playing.")
		return true
	}
}
