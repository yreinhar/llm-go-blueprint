package prompt

import "fmt"

type Prompt interface {
	BuildPromptRequest(userInput, model, task string) (PromptRequest, error)
}

type PromptBuilder struct {
	promptTemplates map[string]PromptTemplate
}

func NewPromptBuilder(files []string) (*PromptBuilder, error) {
	promptTemplates, err := loadPromptTemplates(files)
	if err != nil {
		return &PromptBuilder{}, fmt.Errorf("failed to load prompt templates: %w", err)
	}

	return &PromptBuilder{
		promptTemplates: promptTemplates,
	}, nil
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type PromptRequest struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
}

// BuildPromptRequest builds a prompt request for the given user input and prompt template.
func (pb *PromptBuilder) BuildPromptRequest(userInput, model, task string) (PromptRequest, error) {
	if userInput == "" {
		return PromptRequest{}, fmt.Errorf("user input cannot be empty")
	}

	key := generatePromptKey(model, task)
	developerContent := pb.promptTemplates[key].Roles.Developer.Content

	return PromptRequest{
		Messages: []Message{
			{
				Role:    "developer",
				Content: developerContent,
			},
			{
				Role:    "user",
				Content: userInput,
			},
		},
		Model: model,
	}, nil
}
