package llm

import (
	"bytes"
	"errors"
	"io"
	"net/http"
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
