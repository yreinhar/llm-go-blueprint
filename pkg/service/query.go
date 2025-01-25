package service

import (
	"fmt"

	model "github.com/yreinhar/llm-go-blueprint/pkg/llm/model"
	validation "github.com/yreinhar/llm-go-blueprint/pkg/llm/validation"
)

// QueryService handles requests to LanguageModel.
type QueryService struct {
	LlmModel  model.Llm
	Validator validation.Validation
}

// QueryService creates a new query service for the given large language model.
func NewQueryService(modelName string, schemaPaths []string) (*QueryService, error) {
	llmModel, err := model.GetLlmFactory(modelName)
	if err != nil {
		return nil, fmt.Errorf("failed to get llm model:: %w", err)
	}

	validator, err := validation.NewResponseValidator(schemaPaths)
	if err != nil {
		return nil, fmt.Errorf("failed to create response validator: %w", err)
	}

	return &QueryService{
		LlmModel:  llmModel,
		Validator: validator,
	}, nil
}

// ProcessPrompt processes the input prompt using the specified model and can perform validation
// on the LLM response based a specified output schema.
func (s *QueryService) ProcessPrompt(prompt, responseSchema string) (string, error) {
	// TODO: 1. validate input promp 2. sanitize input prompt 3. call model 4. postprocess repsonse/handle/validate response
	response, err := s.LlmModel.CallModel(prompt)
	if err != nil {
		return "", fmt.Errorf("failed to call model: %w", err)
	}

	if err := s.Validator.Validate(responseSchema, response); err != nil {
		return "", fmt.Errorf("failed to validate response: %w", err)
	}

	return string(response), nil
}
