package validation

// Validation is the common interface for all validators (e.g. response or request).
type Validation interface {
	// Validate validates the data against a given schema. Defining a specific schema allows to handle different task from a llm that produces different outputs.
	Validate(schema string, data []byte) error
}

func NewResponseValidator(schemas []string) (Validation, error) {
	validator, err := NewResponseSchemaValidator(schemas)
	return validator, err
}
