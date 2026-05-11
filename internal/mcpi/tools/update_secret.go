package tools

import (
	"fmt"

	"github.com/infisical/go-sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"mit/platform/internal/mcpi"
)

// RegisterUpdateSecret registers the update_secret tool.
func RegisterUpdateSecret(server *mcpi.Server) {
	server.RegisterTool("update_secret", mcp.Tool{
		Name:        "update_secret",
		Description: "Update an existing secret in Infisical",
		InputSchema: map[string]any{
			"type":       "object",
			"properties": map[string]any{
				"project_id": map[string]any{
					"type":        "string",
					"description": "The project ID where the secret is located",
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
					"description": "The key/name of the secret to update",
				},
				"secret_value": map[string]any{
					"type":        "string",
					"description": "The new value for the secret",
				},
			},
			"required": []any{"project_id", "environment", "secret_key", "secret_value"},
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
		secretValue, _ := args["secret_value"].(string)

		// Update using delete + create pattern (SDK doesn't have a direct update)
		_, err := client.Secrets().Delete(infisical.DeleteSecretOptions{
			ProjectID:   projectID,
			Environment: environment,
			SecretPath:  secretPath,
			SecretKey:   secretKey,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to delete old secret: %w", err)
		}

		_, err = client.Secrets().Create(infisical.CreateSecretOptions{
			ProjectID:   projectID,
			Environment: environment,
			SecretPath:  secretPath,
			SecretKey:   secretKey,
			SecretValue: secretValue,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create updated secret: %w", err)
		}

		return map[string]any{
			"success": true,
			"message": fmt.Sprintf("Secret '%s' updated successfully in project '%s', environment '%s'", secretKey, projectID, environment),
			"project_id":  projectID,
			"environment": environment,
			"secret_path": secretPath,
			"secret_key":  secretKey,
		}, nil
	})
}
