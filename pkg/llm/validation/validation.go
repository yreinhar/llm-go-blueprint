package validation

// Validation is the common interface for all validators (e.g. response or request).
type Validation interface {
	// Validate validates the data against a given schema.
	Validate(schema string, data []byte) error
}

func NewResponseValidator(schemas []string) (Validation, error) {
	validator, err := NewResponseSchemaValidator(schemas)
	return validator, err
}
