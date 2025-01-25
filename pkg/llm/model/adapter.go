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
	CallModel(prompt string) ([]byte, error)
}

type LlamaLocal struct {
	modelName string
	baseURL   string
}

func (m *LlamaLocal) CallModel(prompt string) ([]byte, error) {
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
		return []byte{}, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return []byte{}, fmt.Errorf("creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, errors.New("failed to call Model A API")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
