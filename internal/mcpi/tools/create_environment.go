package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"mit/platform/internal/mcpi"
)

// RegisterCreateEnvironment registers the create_environment tool.
func RegisterCreateEnvironment(server *mcpi.Server) {
	server.RegisterTool("create_environment", mcp.Tool{
		Name:        "create_environment",
		Description: "Create a new environment in Infisical",
		InputSchema: map[string]any{
			"type":       "object",
			"properties": map[string]any{
				"project_id": map[string]any{
					"type":        "string",
					"description": "The project ID where the environment will be created",
				},
				"name": map[string]any{
					"type":        "string",
					"description": "The name of the environment to create (e.g., 'dev', 'staging', 'prod')",
				},
			},
			"required": []any{"project_id", "name"},
		},
	}, func(args map[string]any) (map[string]any, error) {
		// Note: Infisical SDK doesn't have a direct CreateEnvironment method
		// This tool is a placeholder for future API support
		return map[string]any{
			"success": false,
			"message": "Creating environments is not currently supported through the Infisical Go SDK. Please use the Infisical dashboard or CLI to create environments.",
		}, nil
	})
}
