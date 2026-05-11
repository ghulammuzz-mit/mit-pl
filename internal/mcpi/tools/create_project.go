package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"mit/platform/internal/mcpi"
)

// RegisterCreateProject registers the create_project tool.
func RegisterCreateProject(server *mcpi.Server) {
	server.RegisterTool("create_project", mcp.Tool{
		Name:        "create_project",
		Description: "Create a new project in Infisical",
		InputSchema: map[string]any{
			"type":       "object",
			"properties": map[string]any{
				"name": map[string]any{
					"type":        "string",
					"description": "The name of the project to create",
				},
			},
			"required": []any{"name"},
		},
	}, func(args map[string]any) (map[string]any, error) {
		// Note: Infisical SDK doesn't have a direct CreateProject method
		// This tool is a placeholder for future API support
		return map[string]any{
			"success": false,
			"message": "Creating projects is not currently supported through the Infisical Go SDK. Please use the Infisical dashboard or CLI to create projects.",
		}, nil
	})
}
