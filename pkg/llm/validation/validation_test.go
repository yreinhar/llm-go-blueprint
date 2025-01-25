package validation

import (
	"testing"
)

func TestResponseValidatorImplementsInterface(t *testing.T) {
	var _ Validation = (*ResponseSchemaValidator)(nil)
}
