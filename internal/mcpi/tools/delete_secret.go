package tools

import (
	"fmt"

	"github.com/infisical/go-sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"mit/platform/internal/mcpi"
)

// RegisterDeleteSecret registers the delete_secret tool.
func RegisterDeleteSecret(server *mcpi.Server) {
	server.RegisterTool("delete_secret", mcp.Tool{
		Name:        "delete_secret",
		Description: "Delete a secret from Infisical",
		InputSchema: map[string]any{
			"type":       "object",
			"properties": map[string]any{
				"project_id": map[string]any{
					"type":        "string",
					"description": "The project ID where secret is located",
				},
				"environment": map[string]any{
					"type":        "string",
					"description": "The environment (e.g., 'dev', 'staging', 'prod')",
				},
				"secret_path": map[string]any{
					"type":        "string",
					"description": "The folder path where the secret is located (default: '/')",
				},
				"secret_key": map[string]any{
					"type":        "string",
					"description": "The key/name of the secret to delete",
				},
			},
			"required": []any{"project_id", "environment", "secret_key"},
		},
	}, func(args map[string]any) (map[string]any, error) {
		client := server.Client()

		projectID, _ := args["project_id"].(string)
		environment, _ := args["environment"].(string)
		secretPath, _ := args["secret_path"].(string)
		if secretPath == "" {
			secretPath = "/"
		}
		secretKey, _ := args["secret_key"].(string)

		_, err := client.Secrets().Delete(infisical.DeleteSecretOptions{
			ProjectID:   projectID,
			Environment: environment,
			SecretPath:  secretPath,
			SecretKey:   secretKey,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to delete secret: %w", err)
		}

		return map[string]any{
			"success": true,
			"message": fmt.Sprintf("Secret '%s' deleted successfully from project '%s', environment '%s'", secretKey, projectID, environment),
			"project_id":  projectID,
			"environment": environment,
			"secret_path": secretPath,
			"secret_key":  secretKey,
		}, nil
	})
}
