package routes

import (
	"net/http"

	"github.com/yreinhar/llm-go-blueprint/pkg/handlers"
)

func AddRoutes(mux *http.ServeMux) {
	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("/hello", handlers.HandleHelloWorld)
	mux.HandleFunc("/query", handlers.CallModelHandler)
}
