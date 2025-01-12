package app

import (
	"net/http"

	"github.com/yreinhar/llm-go-blueprint/pkg/routes"
)

func NewServer() http.Handler {
	mux := http.NewServeMux()
	routes.AddRoutes(mux)
	var handler http.Handler = mux
	return handler
}
