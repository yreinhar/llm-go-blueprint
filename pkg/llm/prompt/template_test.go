package prompt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadPromptTemplates(t *testing.T) {
	templates, err := loadPromptTemplates([]string{"prompts/promptTemplateDefault.yaml"})
	assert.NoError(t, err)
	assert.NotEmpty(t, templates)

	// Test specific templates
	llamaChatKey := generatePromptKey("llama-3-1b-chat", "chat")
	template, exists := templates[llamaChatKey]
	assert.True(t, exists)
	assert.Equal(t, "llama-3-1b-chat", template.Model)
	assert.Equal(t, "chat", template.Task)
	assert.NotEmpty(t, template.Roles.Developer.Content)
}
