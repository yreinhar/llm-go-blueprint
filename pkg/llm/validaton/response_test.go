package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseValidator(t *testing.T) {
	testCases := []struct {
		name    string
		schemas []string
		length  int
	}{
		{
			name:    "schema: animalResponse",
			schemas: []string{"schemas/animalResponse.cue"},
			length:  1,
		},
		{
			name:    "schema: personResponse",
			schemas: []string{"schemas/personResponse.cue"},
			length:  1,
		},
		{
			name:    "multiple schema: animalResponse, personResponse",
			schemas: []string{"schemas/animalResponse.cue", "schemas/personResponse.cue"},
			length:  2,
		},
	}

	for _, tc := range testCases {
		validator, err := NewResponseValidator(tc.schemas)
		assert.Error(t, err)
		assert.NotNil(t, validator)
		assert.NotEmpty(t, validator.schemas)
		assert.Equal(t, tc.length, len(validator.schemas))
	}
}
