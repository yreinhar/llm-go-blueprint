package prompt

import (
	"embed"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

//go:embed prompts/*.yaml
var schemaFS embed.FS

// PromptTemplate represents the YAML configuration structure
type PromptTemplate struct {
	Model  string       `yaml:"model"`
	Task   string       `yaml:"task"`
	Config PromptConfig `yaml:"config"`
	Roles  Roles        `yaml:"roles"`
}

type PromptConfig struct {
	Temperature float64 `yaml:"temperature"`
}

type Roles struct {
	Developer Role  `yaml:"developer"`           // required
	Assistant *Role `yaml:"assistant,omitempty"` // optional
}

type Role struct {
	Content string `yaml:"content"`
}

// generatePromptKey generates a unique key for a prompt template based on the model and task. This is used to identify the prompt template in the map.
func generatePromptKey(model, task string) string {
	return fmt.Sprintf("%s-%s", model, task)
}

func loadPromptTemplates(promptFiles []string) (map[string]PromptTemplate, error) {
	promptTemplates := make(map[string]PromptTemplate)

	for _, file := range promptFiles {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("reading template file: %w", err)
		}

		var template PromptTemplate
		if err := yaml.Unmarshal(data, &template); err != nil {
			return nil, fmt.Errorf("parsing template yaml: %w", err)
		}

		key := generatePromptKey(template.Model, template.Task)
		log.Debugf("Loaded prompt templates: %v", promptTemplates)

		promptTemplates[key] = template
	}

	if len(promptTemplates) == 0 {
		return nil, fmt.Errorf("no prompt templates found")
	}

	return promptTemplates, nil
}
