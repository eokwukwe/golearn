package main

import (
	"eokwukwe/minigrep/core"
	"fmt"
	"os"
)

// printError prints errors to the standard error stream.
func printError(err error) {
	fmt.Fprintln(os.Stderr, err)
}

// printResults prints the search results to the standard output stream.

func main() {
	config, err := core.BuildConfig(os.Args)
	if err != nil {
		printError(fmt.Errorf("application error: %w", err))
		return
	}

	_, e := core.Search(config)
	if e != nil {
		printError(fmt.Errorf("error during search: %w", e))
		return
	}
}
