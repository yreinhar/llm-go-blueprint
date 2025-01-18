package validation

import (
	"bytes"
	"encoding/json"
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/openapi"
	"github.com/sirupsen/logrus"
)

type SchemaValidator struct {
	cueCtx  *cue.Context
	schemas map[string]cue.Value
}

func NewResponseValidator(schemaFiles []string) (*SchemaValidator, error) {
	cueCtx := cuecontext.New()
	schemas := make(map[string]cue.Value)

	for _, schema := range schemaFiles {
		name, err := getPackageName(schema)
		if err != nil {
			return nil, fmt.Errorf("failed to get schema: %w", err)
		}

		logrus.Debugf("Generating OpenAPI schema for %s", schema)
		cueSchema := cueCtx.CompileString(schema)

		// version := "v1"
		// openAPISchema, err := generateOpenAPISchema(cueSchema, title, version)
		// if err != nil {
		// 	return nil, fmt.Errorf("failed to generate openapi schema: %w", err)
		// }

		schemas[name] = cueSchema
	}

	return &SchemaValidator{
		cueCtx:  cueCtx,
		schemas: schemas,
	}, nil
}

func (v *SchemaValidator) Validate(data string, response []byte) error {
	var jsonData interface{}
	if err := json.Unmarshal(response, &jsonData); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// TODO: validate logic

	return nil
}

func getPackageName(schemaName string) (string, error) {
	instances := load.Instances([]string{schemaName}, nil)
	if len(instances) == 0 {
		return "", fmt.Errorf("no instances found for %s", schemaName)
	}

	fmt.Printf("Package name: %s\n", instances[0].PkgName)

	return "", nil
}

func generateOpenAPISchema(schema cue.Value, title, version string) (string, error) {
	info := struct {
		Title   string `json:"title"`
		Version string `json:"version"`
	}{title, version}

	resolveRefs := &openapi.Config{
		Info:             info,
		ExpandReferences: true,
	}

	logrus.Debugf("Generating OpenAPI schema for %s", title)

	openAPISchema, err := openapi.Gen(schema, resolveRefs)
	if err != nil {
		return "", fmt.Errorf("failed to generate openapi schema: %w", err)
	}

	prettyPrintJSON(openAPISchema)

	return string(openAPISchema), nil
}

func prettyPrintJSON(data []byte) {
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, data, "", "  ")
	fmt.Printf("\n%s\n", prettyJSON.String())
}
