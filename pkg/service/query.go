package service

import model "github.com/yreinhar/llm-go-blueprint/pkg/llm/model"

// QueryService handles requests to LanguageModel.
type QueryService struct {
	llmModel model.Llm
}

// QueryService creates a new query service for the given large language model.
func NewQueryService(modelName string) (*QueryService, error) {
	llmModel, err := model.GetLlmFactory(modelName)
	if err != nil {
		return nil, err
	}

	return &QueryService{
		llmModel: llmModel,
	}, nil
}

// ProcessPrompt processes the input prompt using the specified model.
func (s *QueryService) ProcessPrompt(prompt string) (string, error) {
	return s.llmModel.CallModel(prompt)
}
