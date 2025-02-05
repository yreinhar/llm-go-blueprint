package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockQueryService implements the QueryService interface for testing
type MockQueryService struct {
	mock.Mock
}

func (m *MockQueryService) ProcessPrompt(prompt, schemaType, task string) (string, error) {
	args := m.Called(prompt, schemaType, task)
	return args.String(0), args.Error(1)
}

func TestCreateServer(t *testing.T) {
	server, err := NewServer(
		WithPromptTemplates([]string{"prompts/promptTemplateDefault.yaml"}),
	)
	assert.NoError(t, err)
	assert.NotNil(t, server)
	assert.NotNil(t, server.mux)
	assert.NotNil(t, server.handler)
}

func TestServer(t *testing.T) {
	t.Run("server creation", func(t *testing.T) {
		server, err := NewServer()
		assert.NoError(t, err)
		assert.NotNil(t, server)
		assert.NotNil(t, server.mux)
		assert.NotNil(t, server.handler)
	})

	t.Run("handler registration", func(t *testing.T) {
		server, err := NewServer()
		assert.NoError(t, err)

		handler := server.Handler()
		assert.NotNil(t, handler)

		// Test routes are properly registered
		testCases := []struct {
			name           string
			path           string
			expectedStatus int
		}{
			{
				name:           "root path returns 404",
				path:           "/",
				expectedStatus: http.StatusNotFound,
			},
			{
				name:           "hello path returns 200",
				path:           "/hello",
				expectedStatus: http.StatusOK,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, tc.path, nil)
				rr := httptest.NewRecorder()

				handler.ServeHTTP(rr, req)
				assert.Equal(t, tc.expectedStatus, rr.Code)
			})
		}
	})
}
