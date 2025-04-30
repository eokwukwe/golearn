package tools

import (
	"encoding/json"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/invopop/jsonschema"
)

type ToolDefinition struct {
	Name        string                         `json:"name"`
	Description string                         `json:"description"`
	InputSchema anthropic.ToolInputSchemaParam `json:"input_schema"`
	Function    func(input json.RawMessage) (string, error)
}

// GenerateSchema generates a JSON schema for the specified generic type T.
// It uses the jsonschema Reflector to create a schema definition with 
// the following configurations:
// 
// - Disallows additional properties in the schema.
// - Avoids referencing external schema definitions.
//
// The resulting schema is returned as an anthropic.ToolInputSchemaParam, which
//  includes the schema's properties.
func GenerateSchema[T any]() anthropic.ToolInputSchemaParam {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false, // Disallow additional properties in the schema
		DoNotReference:            true,  // Avoid referencing external schema definitions
	}
	var v T

	schema := reflector.Reflect(v) // Generate the schema for the type T

	return anthropic.ToolInputSchemaParam{
		Properties: schema.Properties, // Return the schema's properties
	}
}
