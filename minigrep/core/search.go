package core

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Search performs a search for the query in the provided file.
// It supports both case-sensitive and case-insensitive
// searches based on the ignoreCase flag.
func Search(c *Config) ([]string, error) {
	var query = c.Query

	if query == "" {
		return nil, errors.New("query cannot be empty")
	}

	if c.IgnoreCase {
		query = strings.ToLower(query)
	}

	// Open the file for reading
	file, err := os.Open(c.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var results []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text() // Read the current line

		if c.IgnoreCase {
			if strings.Contains(strings.ToLower(line), query) {
				results = append(results, line)
			}
		} else {
			if strings.Contains(line, query) {
				results = append(results, line)
			}
		}
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	for _, line := range results {
		fmt.Println(line)
	}

	return results, nil
}
