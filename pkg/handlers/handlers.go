package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type RequestPayload struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type ResponsePayload struct {
	Response string `json:"response"`
}

// QueryService defines the interface for processing model prompts.
// Implementations handle the actual interaction with language models.
type QueryService interface {
	ProcessPrompt(prompt, schemaType, task string) (string, error)
}

// Handler manages HTTP request processing and coordinates with the query service.
// It encapsulates all the dependencies needed for handling HTTP requests.
type Handler struct {
	queryService QueryService
}

const (
	// schema type to validate the llm response against
	schemaTypeToValidateAgainst = "personResponse"
	task                        = "chat"
)

// NewHandler creates a new handler instance with the provided query service.
// It initializes the handler with all required dependencies for processing requests.
func NewHandler(queryService QueryService) *Handler {
	return &Handler{
		queryService: queryService,
	}
}

// CallModelHandler handles the REST API call for calling a model.
func (h *Handler) CallModelHandler(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	payloadString := fmt.Sprintf("\n Model: %s \n Prompt: %s \n", payload.Model, payload.Prompt)

	w.Write([]byte(payloadString))

	response, err := h.queryService.ProcessPrompt(payload.Prompt, schemaTypeToValidateAgainst, task)
	if err != nil {
		http.Error(w, "Failed to process prompt: "+err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse := ResponsePayload{
		Response: response,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonResponse)
}

func (h *Handler) HandleHelloWorld(w http.ResponseWriter, r *http.Request) {
	_ = r.Body
	log.Println("Received a non domain request")
	w.Write([]byte("Hello World"))
}
