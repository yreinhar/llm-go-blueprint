package prompt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPromptBuilder(t *testing.T) {
	pb, err := NewPromptBuilder([]string{"prompts/promptTemplateDefault.yaml"})
	assert.NoError(t, err)
	assert.NotNil(t, pb.promptTemplates)
}

func TestBuildPromptRequest(t *testing.T) {
	pb, err := NewPromptBuilder([]string{"prompts/promptTemplateDefault.yaml"})
	assert.NoError(t, err)
	assert.NotNil(t, pb.promptTemplates)

	userInput := "Hello Model"
	model := "llama-3-1b-chat"
	task := "chat"

	req, err := pb.BuildPromptRequest(userInput, model, task)
	assert.NoError(t, err)
	assert.Equal(t, req.Model, model)
	for _, msg := range req.Messages {
		assert.NotEmpty(t, msg.Role)
		assert.NotEmpty(t, msg.Content)
	}

}
