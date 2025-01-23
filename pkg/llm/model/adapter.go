package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Llm is the common interface for all models.
type Llm interface {
	CallModel(prompt string) (string, error)
}

// Llamadapter interacts with Llama's API.
type LlamaAdapter struct{}

// CallModel sends a POST request to Llama's API.
func (m *LlamaAdapter) CallModel(prompt string) (string, error) {
	url := "https://llama-api.com/v1/call"
	requestBody := []byte(`{"prompt": "` + prompt + `"}`)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to call Model A API")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

type LlamaLocal struct {
	modelName string
	baseURL   string
}

func (m *LlamaLocal) CallModel(prompt string) (string, error) {
	url := fmt.Sprintf("%s/chat/completions", m.baseURL)
	// Request body
	requestBody := map[string]interface{}{
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are a helpful assistant.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"model": m.modelName,
	}

	log.Printf("requestBody: %v", requestBody)

	// Convert to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatalf("failed to marshal JSON: %v", err)
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to call Model A API")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
