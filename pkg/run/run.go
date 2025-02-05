package run

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/yreinhar/llm-go-blueprint/pkg/app"
)

const defaultConfigPath = "files/config.yaml"

// All dependencies as explicit parameter.
func Run(
	ctx context.Context,
	args []string, // For handling command line arguments.
	getenv func(string) string, // For getting environment variables.
	// stdin io.Reader, // For reading input.
	stdout, stderr io.Writer, // For writing output.
) error {
	// Parse command line flags.
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	configPath := flags.String("config", defaultConfigPath, "path to default config file")
	if err := flags.Parse(args[1:]); err != nil {
		return fmt.Errorf("parsing flags: %w", err)
	}

	// Load config with precedence: env vars > config file > defaults.
	config, err := loadConfig(*configPath, getenv)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create a new server instance
	srv, err := app.NewServer()
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: srv.Handler(),
	}

	// Start server in a goroutine
	go func() {
		// Fprintf let me specify where to write and is better for testing and flexibility.
		fmt.Fprintf(stdout, "listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(stderr, "error listening and serving: %s\n", err)
		}
	}()

	// Wait for interruption
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	fmt.Fprintln(stdout, "Shutting down server...")
	if err := httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	fmt.Fprintln(stdout, "Server exiting")
	return nil
}
