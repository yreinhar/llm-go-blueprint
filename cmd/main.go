package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/yreinhar/llm-go-blueprint/pkg/app"
)

func main() {
	// Create a new server instance.
	srv := app.NewServer()
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: srv,
	}

	// Start server in a goroutine in the background.
	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	// Wait for interruption signals while the server runs in the background and gracefully shutdown the server.
	quit := make(chan os.Signal, 1)
	// Tell the signal package to notify our channel when SIGINT or SIGTERM is received
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// Block here until a signal is received
	<-quit

	// Graceful shutdown.
	log.Println("Shutting down server...")
	if err := httpServer.Shutdown(context.Background()); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
