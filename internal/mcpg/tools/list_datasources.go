package tools

import (
	"context"
	"fmt"
	"log/slog"
	"mit/platform/internal/grafana"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ListDataSourcesTool returns the tool definition for list_data_sources.
func ListDataSourcesTool() mcp.Tool {
	return mcp.Tool{
		Name:        "list_data_sources",
		Description: "List all available Grafana data sources. Returns data source UIDs, names, types, and URLs. Useful for discovering available Prometheus data sources for querying.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"type": map[string]any{
					"type":        "string",
					"description": "Filter data sources by type (e.g., 'prometheus')",
				},
			},
		},
	}
}

// ListDataSourcesHandler handles the list_data_sources tool call.
func ListDataSourcesHandler(ctx context.Context, client *grafana.Client, args map[string]any) (map[string]any, error) {
	slog.Info("List data sources request", "args", args)
	slog.Info("DEBUG client", "url", client.URL, "api_key_len", len(client.APIKey))

	// Check if type filter is specified
	var filterType string
	if dsType, ok := args["type"].(string); ok {
		filterType = dsType
	}

	var sources []grafana.DataSource
	var err error

	if filterType == "prometheus" {
		// Find only Prometheus data sources
		sources, err = client.FindPrometheusDataSources(ctx)
	} else {
		// List all data sources
		allSources, err := client.ListDataSources(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list data sources: %w", err)
		}

		// Filter by type if specified
		if filterType != "" {
			sources = make([]grafana.DataSource, 0)
			for _, s := range allSources {
				if s.Type == filterType {
					sources = append(sources, s)
				}
			}
		} else {
			sources = allSources
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list data sources: %w", err)
	}

	slog.Info("Retrieved data sources", "count", len(sources))

	// Format response
	result := make([]map[string]any, len(sources))
	for i, s := range sources {
		result[i] = map[string]any{
			"uid":        s.UID,
			"name":       s.Name,
			"type":       s.Type,
			"url":        s.URL,
			"is_default": s.IsDefault,
		}
	}

	// Find the default data source
	var defaultUID string
	for _, s := range sources {
		if s.IsDefault {
			defaultUID = s.UID
			break
		}
	}

	return map[string]any{
		"data_sources": result,
		"count":        len(result),
		"default_uid":  defaultUID,
	}, nil
}
