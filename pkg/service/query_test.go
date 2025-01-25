package service_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yreinhar/llm-go-blueprint/pkg/llm/prompt"
	"github.com/yreinhar/llm-go-blueprint/pkg/service"
)

type MockLLM struct {
	mock.Mock
}

type MockValidator struct {
	mock.Mock
}

type MockPromptBuilder struct {
	mock.Mock
}

func (m *MockPromptBuilder) BuildPromptRequest(userInput, model, task string) (prompt.PromptRequest, error) {
	args := m.Called(userInput, model, task)
	return args.Get(0).(prompt.PromptRequest), args.Error(1)
}

func (m *MockLLM) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockLLM) CallModel(prompt prompt.PromptRequest) ([]byte, error) {
	args := m.Called(prompt)
	// Get the first argument as []byte directly
	if bytes, ok := args.Get(0).([]byte); ok {
		return bytes, args.Error(1)
	}

	return []byte(args.String(0)), args.Error(1)
}

func (v *MockValidator) Validate(schema string, data []byte) error {
	args := v.Called(schema, data)
	return args.Error(0)
}

func TestNewQueryServiceValid(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	tests := []struct {
		name        string
		modelName   string
		schemaPaths []string
		promptFiles []string
	}{
		{
			name:        "valid model LlamaLocal with empty schema",
			modelName:   "LlamaLocal",
			schemaPaths: []string{},
			promptFiles: []string{},
		},
		{
			name:        "valid model LlamaLocal with schema and prompt files",
			modelName:   "LlamaLocal",
			schemaPaths: []string{"schemas/personResponse.cue"},
			promptFiles: []string{"prompts/promptTemplateDefault.yaml"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := service.NewQueryService(tt.modelName, tt.schemaPaths, tt.promptFiles)
			assert.NoError(t, err)
			assert.NotNil(t, service)
		})
	}
}

func TestNewQueryServiceInvalidModel(t *testing.T) {
	tests := struct {
		name        string
		modelName   string
		schemaPaths []string
		promptFiles []string
	}{
		name:        "model does not exist",
		modelName:   "InvalidModel",
		schemaPaths: []string{},
		promptFiles: []string{},
	}

	t.Run(tests.name, func(t *testing.T) {
		service, err := service.NewQueryService(tests.modelName, tests.schemaPaths, tests.promptFiles)
		assert.Error(t, err)
		assert.Nil(t, service)
	})
}

func TestNewQueryServiceInvalid(t *testing.T) {
	tests := []struct {
		name        string
		modelName   string
		schemaPaths []string
		promptFiles []string
	}{
		{
			name:        "model is valid but schema does not exists",
			modelName:   "LlamaLocal",
			schemaPaths: []string{"schemas/thisFileDoesNotExist.cue"},
			promptFiles: []string{"prompts/promptTemplateDefault.yaml"},
		},
		{
			name:        "model is valid but prompt file does not exist",
			modelName:   "LlamaLocal",
			schemaPaths: []string{"schemas/personResponse.cue"},
			promptFiles: []string{"prompts/thisFileDoesNotExist.yaml"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := service.NewQueryService(tt.modelName, tt.schemaPaths, tt.promptFiles)
			assert.Error(t, err)
			assert.Nil(t, service)
		})
	}
}

func TestQueryServiceProcessPromptSuccess(t *testing.T) {
	testCase := struct {
		name      string
		prompt    string
		modelName string
		mockResp  []byte
	}{

		name:      "successful query",
		prompt:    "What is the capital of France?",
		modelName: "LlamaLocal",
		mockResp:  []byte("The capital of France is Paris."),
	}

	t.Run(testCase.name, func(t *testing.T) {
		// the schema and task itself are not relevant for the test, so mock.Anything is used
		schema := mock.Anything
		task := mock.Anything

		mockPromptBuilder := new(MockPromptBuilder) // Todo
		mockPromptBuilder.On("BuildPromptRequest", testCase.prompt, testCase.modelName, "").Return(prompt.PromptRequest{}, nil)

		mockLLM := new(MockLLM)
		mockLLM.On("CallModel", testCase.prompt).Return(testCase.mockResp, nil)

		mockValidator := new(MockValidator)
		mockValidator.On("Validate", schema, testCase.mockResp).Return(nil)

		// Create service with mock
		service := &service.QueryService{
			LlmModel:      mockLLM,
			Validator:     mockValidator,
			PromptBuilder: mockPromptBuilder,
		}

		got, err := service.ProcessPrompt(testCase.prompt, schema, task)
		assert.NoError(t, err)
		assert.Equal(t, string(testCase.mockResp), got)

		// Verify mock was called as expected
		mockLLM.AssertExpectations(t)
		mockValidator.AssertExpectations(t)
	})
}

func TestQueryServiceProcessPromptCallModelError(t *testing.T) {
	testCase := struct {
		name      string
		prompt    string
		modelName string
		mockResp  []byte
	}{

		name:      "call model error",
		prompt:    "Invalid prompt",
		modelName: "LlamaLocal",
		mockResp:  []byte{},
	}

	t.Run(testCase.name, func(t *testing.T) {
		// the schema and task itself are not relevant for the test, so mock.Anything is used
		schema := mock.Anything
		task := mock.Anything

		mockLLM := new(MockLLM)
		mockLLM.On("CallModel", testCase.prompt).Return(testCase.mockResp, assert.AnError)

		service := &service.QueryService{
			LlmModel:  mockLLM,
			Validator: nil,
		}

		// Model call failed
		_, err := service.ProcessPrompt(testCase.prompt, schema, task)
		assert.Error(t, err)

		// Verify mock was called as expected
		mockLLM.AssertExpectations(t)
	})
}

func TestQueryServiceProcessPromptValidateError(t *testing.T) {
	testCase := struct {
		name      string
		prompt    string
		modelName string
		mockResp  []byte
	}{

		name:      "validate error",
		prompt:    "This is a valid prompt",
		modelName: "LlamaLocal",
		mockResp:  []byte("But this is not a valid response and should fail"),
	}

	t.Run(testCase.name, func(t *testing.T) {
		// the schema and task itself are not relevant for the test, so mock.Anything is used
		schema := mock.Anything
		task := mock.Anything

		mockLLM := new(MockLLM)
		mockLLM.On("CallModel", testCase.prompt).Return(testCase.mockResp, nil)

		mockValidator := new(MockValidator)
		mockValidator.On("Validate", schema, testCase.mockResp).Return(assert.AnError)

		// Create service with mock
		service := &service.QueryService{
			LlmModel:  mockLLM,
			Validator: mockValidator,
		}

		// Model call was successful, but validation failed
		_, err := service.ProcessPrompt(testCase.prompt, schema, task)
		assert.Error(t, err)

		// Verify mock was called as expected
		mockLLM.AssertExpectations(t)
		mockValidator.AssertExpectations(t)
	})
}
