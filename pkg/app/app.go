package app

import (
	"net/http"

	"github.com/yreinhar/llm-go-blueprint/pkg/middleware"
	"github.com/yreinhar/llm-go-blueprint/pkg/routes"
)

func NewServer() http.Handler {
	mux := http.NewServeMux()
	routes.AddRoutes(mux)
	var handler http.Handler = mux
	handler = middleware.LoggingMiddleware(handler) // middlewre wraps the existing handler
	return handler
}
