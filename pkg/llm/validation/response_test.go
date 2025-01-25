package validation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNonExistingSchema(t *testing.T) {
	validator, err := NewResponseSchemaValidator([]string{"schemas/nonExistingSchema.cue"})
	assert.Error(t, err)
	assert.Nil(t, validator)
}

func TestResponseValidatorSchemaLoading(t *testing.T) {
	testCases := []struct {
		name    string
		schemas []string
		length  int
	}{
		{
			name:    "one schema: animalResponse",
			schemas: []string{"schemas/animalResponse.cue"},
			length:  1,
		},
		{
			name:    "one schema: personResponse",
			schemas: []string{"schemas/personResponse.cue"},
			length:  1,
		},
		{
			name:    "multiple schemas: animalResponse, personResponse",
			schemas: []string{"schemas/animalResponse.cue", "schemas/personResponse.cue"},
			length:  2,
		},
	}

	for _, tc := range testCases {
		validator, err := NewResponseSchemaValidator(tc.schemas)
		assert.NoError(t, err)
		assert.NotNil(t, validator)
		assert.NotEmpty(t, validator.schemas)
		assert.Equal(t, tc.length, len(validator.schemas))
	}
}

func TestValidatePersonResponseInvalid(t *testing.T) {
	testCases := []struct {
		name     string
		response []byte
	}{
		{
			name:     "invalid response with unknown properties",
			response: []byte(`{"message": "Hello", "type": "chat", "extra": "extra"}`),
		},
		{
			name:     "invalid response with missing name",
			response: []byte(`{"age": 13}`),
		},
		{
			name:     "invalid response with missing age",
			response: []byte(`{"name": "Hello"}`),
		},
		{
			name:     "invalid response with age out of range",
			response: []byte(`{"name": "Peter", "age": 200}`),
		},
		{
			name:     "invalid response with int as name",
			response: []byte(`{"name": 123, "age": "200"}`),
		},
		{
			name:     "invalid response with age as string",
			response: []byte(`{"name": "Tom", "age": "200"}`),
		},
	}

	responseType := "personResponse"
	schemas := []string{fmt.Sprintf("schemas/%s.cue", responseType)}

	validator, err := NewResponseSchemaValidator(schemas)
	assert.NoError(t, err)
	assert.NotNil(t, validator)
	assert.NotEmpty(t, validator.schemas)

	for _, tc := range testCases {
		err = validator.Validate(responseType, tc.response)
		assert.Error(t, err)
	}
}

func TestValidatePersonResponseValid(t *testing.T) {
	testCases := []struct {
		response []byte
	}{
		{
			response: []byte(`{"name": "Ron", "age": 56}`),
		},
		{
			response: []byte(`{"name": "Harry", "age": 4}`),
		},
		{
			response: []byte(`{"name": "Hermine", "age": 34}`),
		},
		{
			response: []byte(`{"name": "Dobby", "age": 97}`),
		},
		{
			response: []byte(`{"name": "Luna", "age": 22}`),
		},
	}

	responseType := "personResponse"
	schemas := []string{fmt.Sprintf("schemas/%s.cue", responseType)}

	validator, err := NewResponseSchemaValidator(schemas)
	assert.NoError(t, err)
	assert.NotNil(t, validator)
	assert.NotEmpty(t, validator.schemas)

	for _, tc := range testCases {
		err = validator.Validate(responseType, tc.response)
		assert.NoError(t, err)
	}
}

func TestValidateAnimalResponseInvalid(t *testing.T) {
	testCases := []struct {
		name     string
		response []byte
	}{
		{
			name:     "invalid response with unknown properties",
			response: []byte(`{"message": "Hello", "type": "chat", "extra": "extra"}`),
		},
		{
			name:     "invalid response with missing name",
			response: []byte(`{"age": 13}`),
		},
		{
			name:     "invalid response with missing age",
			response: []byte(`{"name": "Dog"}`),
		},
		{
			name:     "invalid response with age out of range",
			response: []byte(`{"name": "Cat", "age": 200}`),
		},
		{
			name:     "invalid response with int as name",
			response: []byte(`{"name": Bee, "age": "200"}`),
		},
		{
			name:     "invalid response with age as string",
			response: []byte(`{"name": "Gopher", "age": "60"}`),
		},
	}

	responseType := "animalResponse"
	schemas := []string{fmt.Sprintf("schemas/%s.cue", responseType)}

	validator, err := NewResponseSchemaValidator(schemas)
	assert.NoError(t, err)
	assert.NotNil(t, validator)
	assert.NotEmpty(t, validator.schemas)

	for _, tc := range testCases {
		err = validator.Validate(responseType, tc.response)
		assert.Error(t, err)
	}
}

func TestValidateAnimalResponseValid(t *testing.T) {
	testCases := []struct {
		response []byte
	}{
		{
			response: []byte(`{"name": "Fox", "age": 12}`),
		},
		{
			response: []byte(`{"name": "Bee", "age": 4}`),
		},
		{
			response: []byte(`{"name": "Eagle", "age": 34}`),
		},
		{
			response: []byte(`{"name": "Badger", "age": 40}`),
		},
		{
			response: []byte(`{"name": "Squirrel", "age": 28}`),
		},
	}

	responseType := "animalResponse"
	schemas := []string{fmt.Sprintf("schemas/%s.cue", responseType)}

	validator, err := NewResponseSchemaValidator(schemas)
	assert.NoError(t, err)
	assert.NotNil(t, validator)
	assert.NotEmpty(t, validator.schemas)

	for _, tc := range testCases {
		err = validator.Validate(responseType, tc.response)
		assert.NoError(t, err)
	}
}

func TestGetPackageNameValidInput(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    string
		wantErr bool
	}{
		{
			name:  "valid package name",
			input: []byte("package personResponse\n\n// A Person Response"),
			want:  "personResponse",
		},
		{
			name:  "package name with underscore",
			input: []byte("package person_response"),
			want:  "person_response",
		},
		{
			name:  "package name with multiple spaces",
			input: []byte("package    personResponse"),
			want:  "personResponse",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkgName, err := getPackageName(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, pkgName)
		})
	}
}

func TestGetPackageNameInvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{

		{
			name:  "missing package declaration",
			input: []byte("// Just a comment\ntype Person struct{}"),
		},
		{
			name:  "empty input",
			input: []byte(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkgName, err := getPackageName(tt.input)
			assert.Error(t, err)
			assert.Empty(t, pkgName)

		})
	}
}
