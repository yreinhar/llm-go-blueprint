package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/yreinhar/llm-go-blueprint/pkg/service"
)

type RequestPayload struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type ResponsePayload struct {
	Response string `json:"response"`
}

// CallModelHandler handles the REST API call for calling a model.
func CallModelHandler(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	payloadString := fmt.Sprintf("\n Model: %s \n Prompt: %s \n", payload.Model, payload.Prompt)

	w.Write([]byte(payloadString))

	service, err := service.NewQueryService(payload.Model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := service.ProcessPrompt(payload.Prompt)
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

func HandleHelloWorld(w http.ResponseWriter, r *http.Request) {
	_ = r.Body
	log.Println("Received a non domain request")
	w.Write([]byte("Hello World"))
}
