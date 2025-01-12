package llm

import "errors"

// GetLlmFactory returns the appropriate Llm implementation.
func GetLlmFactory(modelName string) (Llm, error) {
	switch modelName {
	case "Llama":
		return &LlamaAdapter{}, nil
	default:
		return nil, errors.New("unknown model: " + modelName)
	}
}
