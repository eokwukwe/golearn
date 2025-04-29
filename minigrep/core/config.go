package core

import (
	"errors"
	"os"
)

type Config struct {
	Query      string
	Filename   string
	IgnoreCase bool
}

func BuildConfig(args []string) (*Config, error) {
	if len(args) < 3 {
		return nil, errors.New("not enough arguments. go run . <search word> <file to search> > <output file>")
	}

	ignoreCase := false
	query, filename := args[1], args[2]

	if os.Getenv("IGNORE_CASE") == "true" || os.Getenv("IGNORE_CASE") == "1" {
		ignoreCase = true
	}

	return &Config{
		Query:      query,
		Filename:   filename,
		IgnoreCase: ignoreCase,
	}, nil
}
