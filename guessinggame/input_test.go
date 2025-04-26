package main

import (
	"bufio" // Needed for mocking input
	"fmt"
	"strings"
	"testing" // Needed for capturing output
	// Needed for capturing output
)

// Helper function to create a mock InputHandler with specific input
func mockInputHandler(input string) *InputHandler {
	reader := bufio.NewReader(strings.NewReader(input))
	return &InputHandler{reader: reader}
}

// TestNewInputHandler tests the NewInputHandler constructor
func TestNewInputHandler(t *testing.T) {
	handler := NewInputHandler()
	if handler == nil {
		t.Fatal("NewInputHandler returned nil")
	}
	if handler.reader == nil {
		t.Error("NewInputHandler did not initialize the reader")
	}
	// More detailed checks on the reader might be complex and less valuable
}

// TestReadInput tests the ReadInput method with various inputs
func TestReadInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected InputResult
	}{
		{
			name:  "ValidNumber",
			input: "42\n",
			expected: InputResult{
				value:   42,
				isEmpty: false,
				isQuit:  false,
				err:     nil,
			},
		},
		{
			name:  "ValidNumberWithWhitespace",
			input: "  99 \n",
			expected: InputResult{
				value:   99,
				isEmpty: false,
				isQuit:  false,
				err:     nil,
			},
		},
		{
			name:  "EmptyInput",
			input: "\n",
			expected: InputResult{
				value:   0,
				isEmpty: true,
				isQuit:  false,
				err:     nil,
			},
		},
		{
			name:  "QuitCommandLower",
			input: "quit\n",
			expected: InputResult{
				value:   0,
				isEmpty: false,
				isQuit:  true,
				err:     nil,
			},
		},
		{
			name:  "QuitCommandUpper",
			input: "QUIT\n",
			expected: InputResult{
				value:   0,
				isEmpty: false,
				isQuit:  true,
				err:     nil,
			},
		},
		{
			name:  "InvalidNumber",
			input: "abc\n",
			expected: InputResult{
				value:   0,
				isEmpty: false,
				isQuit:  false,
				// We expect an error, but don't check the exact error message
				// as it might change with Go versions or fmt.Errorf wrapping.
				// We just check that err is not nil.
				err: fmt.Errorf("invalid number format: strconv.Atoi: parsing \"abc\": invalid syntax"), // Example error structure
			},
		},
		{
			name:  "InvalidNumberWithWhitespace",
			input: "  xyz \n",
			expected: InputResult{
				value:   0,
				isEmpty: false,
				isQuit:  false,
				err:     fmt.Errorf("invalid number format: strconv.Atoi: parsing \"xyz\": invalid syntax"), // Example error structure
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := mockInputHandler(tt.input)
			result := handler.ReadInput()

			if result.value != tt.expected.value {
				t.Errorf("Expected value %d, got %d", tt.expected.value, result.value)
			}
			if result.isEmpty != tt.expected.isEmpty {
				t.Errorf("Expected isEmpty %t, got %t", tt.expected.isEmpty, result.isEmpty)
			}
			if result.isQuit != tt.expected.isQuit {
				t.Errorf("Expected isQuit %t, got %t", tt.expected.isQuit, result.isQuit)
			}

			// Check error presence, not the exact error value
			if (result.err != nil) != (tt.expected.err != nil) {
				t.Errorf("Expected error presence %t, got %t (err: %v)", tt.expected.err != nil, result.err != nil, result.err)
			}
			// Optional: If you need to check the error message content
			if result.err != nil && tt.expected.err != nil && !strings.Contains(result.err.Error(), tt.expected.err.Error()) {
				t.Errorf("Expected error message to contain '%v', got '%v'", tt.expected.err, result.err)
			}
		})
	}
}
