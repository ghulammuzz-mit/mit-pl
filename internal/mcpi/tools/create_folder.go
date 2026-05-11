package tools

import (
	"fmt"

	"github.com/infisical/go-sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"mit/platform/internal/mcpi"
)

// RegisterCreateFolder registers the create_folder tool.
func RegisterCreateFolder(server *mcpi.Server) {
	server.RegisterTool("create_folder", mcp.Tool{
		Name:        "create_folder",
		Description: "Create a new folder in Infisical",
		InputSchema: map[string]any{
			"type":       "object",
			"properties": map[string]any{
				"project_id": map[string]any{
					"type":        "string",
					"description": "The project ID where the folder will be created",
				},
				"environment": map[string]any{
					"type":        "string",
					"description": "The environment where the folder will be created",
				},
				"name": map[string]any{
					"type":        "string",
					"description": "The name of the folder to create",
				},
			},
			"required": []any{"project_id", "environment", "name"},
		},
	}, func(args map[string]any) (map[string]any, error) {
		client := server.Client()

		projectID, _ := args["project_id"].(string)
		environment, _ := args["environment"].(string)
		name, _ := args["name"].(string)

		_, err := client.Folders().Create(infisical.CreateFolderOptions{
			ProjectID:   projectID,
			Environment: environment,
			Name:        name,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create folder: %w", err)
		}

		return map[string]any{
			"success": true,
			"message": fmt.Sprintf("Folder '%s' created successfully in project '%s', environment '%s'", name, projectID, environment),
			"project_id":  projectID,
			"environment": environment,
			"folder_name": name,
		}, nil
	})
}
