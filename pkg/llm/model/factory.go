package model

import "errors"

// GetLlmFactory returns the appropriate Llm implementation.
func GetLlmFactory(modelName string) (Llm, error) {
	switch modelName {
	case "LlamaLocal":
		return &LlamaLocal{
			modelName: "llama-3-1b-chat",
			baseURL:   "http://localhost:8080/v1",
		}, nil
	default:
		return nil, errors.New("unknown model: " + modelName)
	}
}
