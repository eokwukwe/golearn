package core

import (
	"os"
	"testing"
)

func TestSearch_CaseSensitive(t *testing.T) {
	config := &Config{
		Query:      "line",
		Filename:   "testfile.txt",
		IgnoreCase: false,
	}

	// Create a temporary test file
	fileContent := "This is the first line\nAnother line with the word\nYet another line"
	tempFile, err := os.Create(config.Filename)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.WriteString(fileContent)
	tempFile.Close()

	results, err := Search(config)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := []string{
		"This is the first line",
		"Another line with the word",
		"Yet another line",
	}

	if len(results) != len(expected) {
		t.Errorf("Expected %d results, got %d", len(expected), len(results))
	}

	for i, result := range results {
		if result != expected[i] {
			t.Errorf("Expected result %q, got %q", expected[i], result)
		}
	}
}

func TestSearch_CaseInsensitive(t *testing.T) {
	config := &Config{
		Query:      "line",
		Filename:   "testfile.txt",
		IgnoreCase: true,
	}

	// Create a temporary test file
	fileContent := "This is the first line\nAnother Line with the word\nYet another line"
	tempFile, err := os.Create(config.Filename)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.WriteString(fileContent)
	tempFile.Close()

	results, err := Search(config)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := []string{
		"This is the first line",
		"Another Line with the word",
		"Yet another line",
	}

	if len(results) != len(expected) {
		t.Errorf("Expected %d results, got %d", len(expected), len(results))
	}

	for i, result := range results {
		if result != expected[i] {
			t.Errorf("Expected result %q, got %q", expected[i], result)
		}
	}
}
