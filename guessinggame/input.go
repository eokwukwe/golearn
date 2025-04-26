package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// InputResult represents the result of processing user input
type InputResult struct {
	value   int
	isEmpty bool
	isQuit  bool
	err     error
}

// InputHandler handles user input
type InputHandler struct {
	reader *bufio.Reader
}

// NewInputHandler creates a new input handler
func NewInputHandler() *InputHandler {
	return &InputHandler{
		reader: bufio.NewReader(os.Stdin),
	}
}

// ReadInput reads and processes user input
func (h *InputHandler) ReadInput() InputResult {
	input, err := h.reader.ReadString('\n')

	if err != nil {
		return InputResult{err: fmt.Errorf("error reading input: %w", err)}
	}

	input = strings.TrimSpace(input)

	// Check for empty input
	if len(input) == 0 {
		return InputResult{isEmpty: true}
	}

	// Check for quit command
	if strings.ToLower(input) == "quit" {
		return InputResult{isQuit: true}
	}

	// Try to convert to number
	guess, err := strconv.Atoi(input)

	if err != nil {
		return InputResult{err: fmt.Errorf("invalid number format: %w", err)}
	}

	// Return successful result
	return InputResult{value: guess}
}
