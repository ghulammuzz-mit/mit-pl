package tools

import (
	"fmt"

	"github.com/infisical/go-sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"mit/platform/internal/mcpi"
)

// RegisterGetSecret registers the get_secret tool.
func RegisterGetSecret(server *mcpi.Server) {
	server.RegisterTool("get_secret", mcp.Tool{
		Name:        "get_secret",
		Description: "Get a single secret from Infisical",
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
					"description": "The key/name of the secret to retrieve",
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

		// List all secrets and find the one we need
		secrets, err := client.Secrets().List(infisical.ListSecretsOptions{
			ProjectID:          projectID,
			Environment:         environment,
			SecretPath:         secretPath,
			AttachToProcessEnv: false,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list secrets: %w", err)
		}

		// Find the secret with matching key
		for _, secret := range secrets {
			if secret.SecretKey == secretKey {
				return map[string]any{
					"success": true,
					"project_id":  projectID,
					"environment": environment,
					"secret_path": secretPath,
					"secret_key":  secretKey,
					"secret_value": secret.SecretValue,
					"version": secret.Version,
				}, nil
			}
		}

		return nil, fmt.Errorf("secret '%s' not found in project '%s', environment '%s'", secretKey, projectID, environment)
	})
}
