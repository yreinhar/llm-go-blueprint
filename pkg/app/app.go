package app

import (
	"fmt"
	"net/http"

	"github.com/yreinhar/llm-go-blueprint/pkg/handlers"
	"github.com/yreinhar/llm-go-blueprint/pkg/middleware"
	"github.com/yreinhar/llm-go-blueprint/pkg/routes"
	"github.com/yreinhar/llm-go-blueprint/pkg/service"
)

// Server represents the HTTP server with its dependencies.
// It holds the handler for processing requests and the mux for routing.
type Server struct {
	handler *handlers.Handler
	mux     *http.ServeMux
}

// ServerOption represents a server configuration option
type ServerOption func(*serverConfig)

type serverConfig struct {
	model           string
	responseSchemas []string
	promptTemplates []string
}

// WithModel sets the model for the query service
func WithModel(model string) ServerOption {
	return func(c *serverConfig) {
		c.model = model
	}
}

// WithResponseSchemas sets the response schemas for validation
func WithResponseSchemas(schemas []string) ServerOption {
	return func(c *serverConfig) {
		c.responseSchemas = schemas
	}
}

// WithPromptTemplates sets the prompt templates
func WithPromptTemplates(templates []string) ServerOption {
	return func(c *serverConfig) {
		c.promptTemplates = templates
	}
}

const (
	personResponseSchema = "schemas/personResponse.cue"
	promptTestTemplate   = "prompts/promptTemplateDefault.yaml"
)

// NewServer creates a new server instance with all required dependencies.
// It initializes the query service and handler, returning an error if initialization fails.
func NewServer(opts ...ServerOption) (*Server, error) {
	// Default configuration
	cfg := &serverConfig{
		model:           "LlamaLocal",
		responseSchemas: []string{personResponseSchema},
		promptTemplates: []string{promptTestTemplate},
	}

	// Apply options
	for _, opt := range opts {
		opt(cfg)
	}

	// Create query service with configuration
	queryService, err := service.NewQueryService(
		cfg.model,
		cfg.responseSchemas,
		cfg.promptTemplates,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create query service: %w", err)
	}

	return &Server{
		handler: handlers.NewHandler(queryService),
		mux:     http.NewServeMux(),
	}, nil
}

// Handler returns the configured http.Handler with all routes and middleware applied.
// It sets up the routes and wraps the handler with logging middleware.
func (s *Server) Handler() http.Handler {
	routes.AddRoutes(s.mux, s.handler)
	// middlewre wraps the existing handler
	return middleware.LoggingMiddleware(s.mux)
}
