package tools

import (
	"fmt"

	"github.com/infisical/go-sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"mit/platform/internal/mcpi"
)

// RegisterListSecrets registers the list_secrets tool.
func RegisterListSecrets(server *mcpi.Server) {
	server.RegisterTool("list_secrets", mcp.Tool{
		Name:        "list_secrets",
		Description: "List all secrets in a project/environment",
		InputSchema: map[string]any{
			"type":       "object",
			"properties": map[string]any{
				"project_id": map[string]any{
					"type":        "string",
					"description": "The project ID to list secrets from",
				},
				"environment": map[string]any{
					"type":        "string",
					"description": "The environment (e.g., 'dev', 'staging', 'prod')",
				},
				"secret_path": map[string]any{
					"type":        "string",
					"description": "The folder path to list secrets from (default: '/')",
				},
			},
			"required": []any{"project_id", "environment"},
		},
	}, func(args map[string]any) (map[string]any, error) {
		client := server.Client()

		projectID, _ := args["project_id"].(string)
		environment, _ := args["environment"].(string)
		secretPath, _ := args["secret_path"].(string)
		if secretPath == "" {
			secretPath = "/"
		}

		secrets, err := client.Secrets().List(infisical.ListSecretsOptions{
			ProjectID:          projectID,
			Environment:         environment,
			SecretPath:         secretPath,
			AttachToProcessEnv: false,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list secrets: %w", err)
		}

		// Convert secrets to a simpler format for JSON response
		secretList := make([]map[string]any, 0, len(secrets))
		for _, s := range secrets {
			secretList = append(secretList, map[string]any{
				"key":   s.SecretKey,
				"path":  s.SecretPath,
				"version": s.Version,
			})
		}

		return map[string]any{
			"success": true,
			"project_id":  projectID,
			"environment": environment,
			"secret_path": secretPath,
			"count":       len(secretList),
			"secrets":    secretList,
		}, nil
	})
}
