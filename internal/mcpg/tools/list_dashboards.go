package tools

import (
	"context"
	"fmt"
	"log/slog"
	"mit/platform/internal/grafana"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ListDashboardsTool returns the tool definition for list_dashboards.
func ListDashboardsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "list_dashboards",
		Description: "List available Grafana dashboards with optional search and filter capabilities. Returns dashboard UIDs, titles, folders, and tags.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"query": map[string]any{
					"type":        "string",
					"description": "Search query string to filter dashboards",
				},
				"tag": map[string]any{
					"type":        "string",
					"description": "Filter dashboards by tag",
				},
				"limit": map[string]any{
					"type":        "number",
					"description": "Maximum number of dashboards to return",
				},
			},
		},
	}
}

// ListDashboardsHandler handles the list_dashboards tool call.
func ListDashboardsHandler(ctx context.Context, client *grafana.Client, args map[string]any) (map[string]any, error) {
	slog.Info("List dashboards request", "args", args)

	// Parse parameters
	params := &grafana.SearchDashboardsParams{}

	if query, ok := args["query"].(string); ok {
		params.Query = query
	}

	if tag, ok := args["tag"].(string); ok {
		params.Tags = []string{tag}
	}

	if limit, ok := args["limit"].(float64); ok {
		params.Limit = int(limit)
	}

	// List dashboards
	dashboards, err := client.SearchDashboards(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list dashboards: %w", err)
	}

	slog.Info("Retrieved dashboards", "count", len(dashboards))

	// Format response
	result := make([]map[string]any, len(dashboards))
	for i, d := range dashboards {
		result[i] = map[string]any{
			"uid":     d.UID,
			"title":   d.Title,
			"folder":  d.Folder,
			"tags":    d.Tags,
		}
	}

	return map[string]any{
		"dashboards": result,
		"count":      len(result),
	}, nil
}
