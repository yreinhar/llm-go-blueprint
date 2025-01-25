package validation

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"regexp"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/encoding/openapi"
	"github.com/getkin/kin-openapi/openapi3"
	log "github.com/sirupsen/logrus"
)

type ResponseSchemaValidator struct {
	schemas map[string]*openapi3.Schema
}

//go:embed schemas/*.cue
var schemaFS embed.FS

const version = "v1"

// NewResponseValidator implements the Validation interface and loads all schema files and generates an OpenAPI schema for each that can be used to validate a response. One validator could validate multiple schemas.
func NewResponseSchemaValidator(schemaFiles []string) (*ResponseSchemaValidator, error) {
	cueCtx := cuecontext.New()
	schemas := make(map[string]*openapi3.Schema)

	log.Debugf("Loading schemas: %s", schemaFiles)

	for _, schema := range schemaFiles {
		schemaData, err := schemaFS.ReadFile(schema)
		if err != nil {
			return nil, fmt.Errorf("reading schema file: %w", err)
		}

		name, err := getPackageName(schemaData)
		if err != nil {
			return nil, fmt.Errorf("failed to get package name: %w", err)
		}

		openAPISchema, err := generateOpenAPISchema(cueCtx, schemaData, name, version)
		if err != nil {
			return nil, fmt.Errorf("failed to generate openapi schema: %w", err)
		}

		loader := openapi3.NewLoader()
		doc, err := loader.LoadFromData(openAPISchema)
		if err != nil {
			return nil, fmt.Errorf("loading schema: %w", err)
		}

		if doc.Components.Schemas[name] == nil {
			return nil, fmt.Errorf("schema not found: %s", name)
		}

		setStrictSchemaValidationRules(doc.Components.Schemas[name].Value)
		// Store schema
		schemas[name] = doc.Components.Schemas[name].Value
	}

	return &ResponseSchemaValidator{
		schemas: schemas,
	}, nil
}

func (v *ResponseSchemaValidator) Validate(schema string, data []byte) error {
	// Get the schema from the map
	cueSchema, exists := v.schemas[schema]
	if !exists {
		return fmt.Errorf("validation failed schema not found: %s", schema)
	}

	var jsonData interface{}
	// Parse JSON
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	log.Debugf("Validating against schema: %s\n", schema)

	// Validate against schema
	err := cueSchema.VisitJSON(jsonData)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

// getPackageName gets the package name from the cue schema file.
func getPackageName(data []byte) (string, error) {
	// Parse package declaration using regex
	re := regexp.MustCompile(`package\s+(\w+)`)

	matches := re.FindSubmatch(data)
	if len(matches) < 2 {
		return "", fmt.Errorf("no package name found in schema")
	}

	pkgName := string(matches[1])

	log.Debugf("Found package name: %s", pkgName)

	return pkgName, nil
}

// processSchema loads and compiles the cue schema file.
func processSchema(cueCtx *cue.Context, cueData []byte) (cue.Value, error) {
	// Compile schema
	cueValue := cueCtx.CompileString(string(cueData))
	if cueValue.Err() != nil {
		return cue.Value{}, fmt.Errorf("compiling CUE schema failed: %w", cueValue.Err())
	}

	return cueValue, nil
}

// setStrictSchemaValidationRules sets the validation rules for the schema.
func setStrictSchemaValidationRules(schema *openapi3.Schema) {
	// Disallowing additional properties
	schema.AdditionalProperties = openapi3.AdditionalProperties{
		Has: openapi3.BoolPtr(false),
	}
	// Making all properties required
	required := make([]string, 0)
	for propName := range schema.Properties {
		required = append(required, propName)
	}
	schema.Required = required
}

func generateOpenAPISchema(cueCtx *cue.Context, cueData []byte, title, version string) ([]byte, error) {
	cueSchema := cueCtx.CompileString(string(cueData))
	if cueSchema.Err() != nil {
		return nil, fmt.Errorf("compiling CUE schema failed: %w", cueSchema.Err())
	}

	info := struct {
		Title   string `json:"title"`
		Version string `json:"version"`
	}{title, version}

	resolveRefs := &openapi.Config{
		Info:             info,
		ExpandReferences: true,
	}

	log.Debugf("Generating OpenAPI schema for %s", title)

	openAPISchema, err := openapi.Gen(cueSchema, resolveRefs)
	if err != nil {
		return nil, fmt.Errorf("failed to generate openapi schema: %w", err)
	}

	prettyPrintJSON(openAPISchema)

	return openAPISchema, nil
}

func prettyPrintJSON(data []byte) {
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, data, "", "  ")
	log.Debugf("Generated OpenAPI schema:\n%s\n", prettyJSON.String())
}
