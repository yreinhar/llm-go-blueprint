package main

import (
	"context"
	"fmt"
	"os"

	"github.com/yreinhar/llm-go-blueprint/pkg/run"
)

// main handles only process setup and error reporting.
func main() {
	if err := run.Run(
		context.Background(),
		os.Args,
		os.Getenv,
		// os.Stdin,
		os.Stdout,
		os.Stderr,
	); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
