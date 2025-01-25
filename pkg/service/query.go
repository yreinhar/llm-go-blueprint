package service

import (
	"fmt"

	model "github.com/yreinhar/llm-go-blueprint/pkg/llm/model"
	"github.com/yreinhar/llm-go-blueprint/pkg/llm/prompt"
	validation "github.com/yreinhar/llm-go-blueprint/pkg/llm/validation"
)

// QueryService handles requests to LanguageModel.
type QueryService struct {
	LlmModel      model.Llm
	Validator     validation.Validation
	PromptBuilder prompt.Prompt
}

// QueryService creates a new query service for the given large language model.
func NewQueryService(modelName string, schemaPaths []string, promptFiles []string) (*QueryService, error) {
	llmModel, err := model.GetLlmFactory(modelName)
	if err != nil {
		return nil, fmt.Errorf("failed to get llm model:: %w", err)
	}

	promptBuilder, err := prompt.NewPromptBuilder(promptFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to create prompt builder: %w", err)
	}

	validator, err := validation.NewResponseSchemaValidator(schemaPaths)
	if err != nil {
		return nil, fmt.Errorf("failed to create response validator: %w", err)
	}

	return &QueryService{
		LlmModel:      llmModel,
		Validator:     validator,
		PromptBuilder: promptBuilder,
	}, nil
}

// ProcessPrompt processes the input prompt using the specified model and can perform validation
// on the LLM response based a specified output schema.
func (s *QueryService) ProcessPrompt(prompt, responseSchema, task string) (string, error) {
	// TODO: 1. validate input promp 2. sanitize input prompt 3. call model 4. postprocess repsonse/handle/validate response
	request, err := s.PromptBuilder.BuildPromptRequest(prompt, s.LlmModel.Name(), task)
	if err != nil {
		return "", fmt.Errorf("failed to build prompt request: %w", err)
	}

	response, err := s.LlmModel.CallModel(request)
	if err != nil {
		return "", fmt.Errorf("failed to call model: %w", err)
	}

	if err := s.Validator.Validate(responseSchema, response); err != nil {
		return "", fmt.Errorf("failed to validate response: %w", err)
	}

	return string(response), nil
}
