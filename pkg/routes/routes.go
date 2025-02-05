package routes

import (
	"net/http"

	"github.com/yreinhar/llm-go-blueprint/pkg/handlers"
)

// AddRoutes configures all HTTP routes for the app.
func AddRoutes(mux *http.ServeMux, h *handlers.Handler) {
	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("/hello", h.HandleHelloWorld)
	mux.HandleFunc("/query", h.CallModelHandler)
}
