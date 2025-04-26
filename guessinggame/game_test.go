package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

// TestNewGame tests the NewGame constructor
func TestNewGame(t *testing.T) {
	min := 1
	max := 100
	game := NewGame(min, max)

	// Check if the game struct is created
	if game == nil {
		t.Fatal("NewGame returned nil")
	}

	// Check if the range is set correctly
	if game.minRange != min {
		t.Errorf("Expected minRange %d, got %d", min, game.minRange)
	}
	if game.maxRange != max {
		t.Errorf("Expected maxRange %d, got %d", max, game.maxRange)
	}

	// Check if the secret number is within the specified range
	// Note: We can't test for a specific secret number due to randomness,
	// but we can check if it's within the bounds.
	if game.secretNumber < min || game.secretNumber > max {
		t.Errorf("Secret number %d is outside the expected range [%d, %d]", game.secretNumber, min, max)
	}
}

// TestIsValidGuess tests the isValidGuess method
func TestIsValidGuess(t *testing.T) {
	game := &Game{minRange: 10, maxRange: 20} // Create a dummy game for testing range

	tests := []struct {
		guess    int
		expected bool
	}{
		{10, true},  // Lower bound
		{15, true},  // Within range
		{20, true},  // Upper bound
		{9, false},  // Below range
		{21, false}, // Above range
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Guess_%d", tt.guess), func(t *testing.T) {
			result := game.isValidGuess(tt.guess)
			if result != tt.expected {
				t.Errorf("For guess %d, expected isValidGuess to be %t, got %t", tt.guess, tt.expected, result)
			}
		})
	}
}

// TestCheckGuess tests the checkGuess method
func TestCheckGuess(t *testing.T) {
	// Create a dummy game with a known secret number for testing
	secret := 50
	game := &Game{secretNumber: secret}

	tests := []struct {
		guess    int
		expected bool   // Expected return value of checkGuess
		output   string // Expected output string (for basic check)
	}{
		{40, false, "Too low!"},        // Guess too low
		{60, false, "Too high!"},       // Guess too high
		{50, true, "Congratulations!"}, // Correct guess
	}

	// Redirect standard output to capture fmt.Println
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Guess_%d", tt.guess), func(t *testing.T) {
			// Call the method
			result := game.checkGuess(tt.guess)

			// Check the boolean result
			if result != tt.expected {
				t.Errorf("For guess %d, expected checkGuess to return %t, got %t", tt.guess, tt.expected, result)
			}

			// Capture and check the output (basic check for substring)
			w.Close() // Close the writer to flush output
			out, _ := io.ReadAll(r)
			os.Stdout = oldStdout // Restore original stdout
			outputString := string(out)

			if !strings.Contains(outputString, tt.output) {
				t.Errorf("For guess %d, expected output to contain '%s', got '%s'", tt.guess, tt.output, outputString)
			}

			// Reset the pipe for the next test iteration
			r, w, _ = os.Pipe()
			os.Stdout = w
		})
	}

	// Restore original stdout one last time after the loop
	w.Close() // Close the writer from the last iteration
	os.Stdout = oldStdout
}
