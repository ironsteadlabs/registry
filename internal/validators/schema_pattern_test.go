package validators_test

import (
	"encoding/json"
	"testing"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiv0 "github.com/modelcontextprotocol/registry/pkg/api/v0"
	"github.com/modelcontextprotocol/registry/pkg/model"
)

// TestSchema_URLPatternValidation tests that the JSON schema properly validates URL patterns
// in both StreamableHttpTransport and SseTransport.
//
// The schema should:
// 1. Require URL-like values (pattern: "^https?://[^\\s]+$") to prevent free-form strings
// 2. Allow template variables like {tenant_id} in URLs (which format: "uri" would reject)
func TestSchema_URLPatternValidation(t *testing.T) {
	// Compile the schema
	schemaPath := "../../docs/reference/server-json/server.schema.json"
	compiler := jsonschema.NewCompiler()
	compiler.Draft = jsonschema.Draft7
	schema, err := compiler.Compile(schemaPath)
	require.NoError(t, err, "Failed to compile server.schema.json")

	tests := []struct {
		name        string
		serverJSON  apiv0.ServerJSON
		shouldPass  bool
		description string
	}{
		{
			name: "StreamableHttp with valid static URL should pass",
			serverJSON: apiv0.ServerJSON{
				Schema:      model.CurrentSchemaURL,
				Name:        "com.example/test",
				Description: "Test",
				Version:     "1.0.0",
				Repository: model.Repository{
					URL:    "https://github.com/example/test",
					Source: "github",
				},
				Remotes: []model.Transport{
					{
						Type: "streamable-http",
						URL:  "https://example.com/mcp",
					},
				},
			},
			shouldPass:  true,
			description: "Valid static HTTPS URL should always pass",
		},
		{
			name: "StreamableHttp with template variable should pass",
			serverJSON: apiv0.ServerJSON{
				Schema:      model.CurrentSchemaURL,
				Name:        "com.example/test",
				Description: "Test",
				Version:     "1.0.0",
				Repository: model.Repository{
					URL:    "https://github.com/example/test",
					Source: "github",
				},
				Remotes: []model.Transport{
					{
						Type: "streamable-http",
						URL:  "https://example.com/mcp/{tenant_id}",
						Variables: map[string]model.Input{
							"tenant_id": {
								Description: "Tenant ID",
								IsRequired:  true,
							},
						},
					},
				},
			},
			shouldPass:  true,
			description: "Template variables in URL should be allowed (format: uri would reject this)",
		},
		{
			name: "SSE with template variable should pass",
			serverJSON: apiv0.ServerJSON{
				Schema:      model.CurrentSchemaURL,
				Name:        "com.example/test",
				Description: "Test",
				Version:     "1.0.0",
				Repository: model.Repository{
					URL:    "https://github.com/example/test",
					Source: "github",
				},
				Remotes: []model.Transport{
					{
						Type: "sse",
						URL:  "https://example.com/sse/{tenant_id}",
						Variables: map[string]model.Input{
							"tenant_id": {
								Description: "Tenant ID",
								IsRequired:  true,
							},
						},
					},
				},
			},
			shouldPass:  true,
			description: "SSE with template variable should be allowed (current format: uri rejects this)",
		},
		{
			name: "StreamableHttp with non-URL string should fail",
			serverJSON: apiv0.ServerJSON{
				Schema:      model.CurrentSchemaURL,
				Name:        "com.example/test",
				Description: "Test",
				Version:     "1.0.0",
				Repository: model.Repository{
					URL:    "https://github.com/example/test",
					Source: "github",
				},
				Remotes: []model.Transport{
					{
						Type: "streamable-http",
						URL:  "not a url at all",
					},
				},
			},
			shouldPass:  false,
			description: "Free-form string without http(s):// should be rejected by pattern",
		},
		{
			name: "SSE with non-URL string should fail",
			serverJSON: apiv0.ServerJSON{
				Schema:      model.CurrentSchemaURL,
				Name:        "com.example/test",
				Description: "Test",
				Version:     "1.0.0",
				Repository: model.Repository{
					URL:    "https://github.com/example/test",
					Source: "github",
				},
				Remotes: []model.Transport{
					{
						Type: "sse",
						URL:  "just some text",
					},
				},
			},
			shouldPass:  false,
			description: "Free-form string without http(s):// should be rejected by pattern",
		},
		{
			name: "StreamableHttp with URL containing spaces should fail",
			serverJSON: apiv0.ServerJSON{
				Schema:      model.CurrentSchemaURL,
				Name:        "com.example/test",
				Description: "Test",
				Version:     "1.0.0",
				Repository: model.Repository{
					URL:    "https://github.com/example/test",
					Source: "github",
				},
				Remotes: []model.Transport{
					{
						Type: "streamable-http",
						URL:  "https://example.com/mcp with spaces",
					},
				},
			},
			shouldPass:  false,
			description: "URL with spaces should be rejected by pattern",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON and back to get the raw data structure
			jsonBytes, err := json.Marshal(tt.serverJSON)
			require.NoError(t, err, "Failed to marshal server JSON")

			var jsonData interface{}
			err = json.Unmarshal(jsonBytes, &jsonData)
			require.NoError(t, err, "Failed to unmarshal JSON data")

			// Validate against schema
			err = schema.Validate(jsonData)

			if tt.shouldPass {
				assert.NoError(t, err, "%s: %v", tt.description, err)
			} else {
				assert.Error(t, err, "%s: expected schema validation to fail", tt.description)
			}
		})
	}
}
