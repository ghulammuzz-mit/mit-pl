package tools

import (
	"context"
	"fmt"
	"log/slog"
	"mit/platform/internal/grafana"
	"mit/platform/internal/mcpg"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// QueryMetricsTool returns the tool definition for query_metrics.
func QueryMetricsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "query_metrics",
		Description: "Query Prometheus metrics from Grafana. Supports natural language queries like 'p95 of service X in namespace Y for the last week'. Returns time series data with timestamps and values.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"query": map[string]any{
					"type":        "string",
					"description": "The PromQL query or metric name (e.g., 'http_request_duration_seconds', 'rate(http_requests_total[5m])')",
				},
				"percentile": map[string]any{
					"type":        "number",
					"description": "Percentile for calculations (e.g., 95 for p95, 99 for p99). Only used with 'percentile' aggregation.",
				},
				"timeRange": map[string]any{
					"type":        "string",
					"description": "Time range: '1h', '6h', '24h', '7d', '30d'. Default: '24h'",
					"enum":        []any{"1h", "6h", "24h", "7d", "30d"},
				},
				"namespace": map[string]any{
					"type":        "string",
					"description": "Kubernetes namespace to filter on",
				},
				"service": map[string]any{
					"type":        "string",
					"description": "Service name to filter on",
				},
				"dataSource": map[string]any{
					"type":        "string",
					"description": "Grafana data source UID. If not specified, uses the default Prometheus data source",
				},
				"aggregation": map[string]any{
					"type":        "string",
					"description": "Aggregation type to apply",
					"enum":        []any{"avg", "sum", "max", "min", "count", "percentile"},
				},
			},
			"required": []any{"query"},
		},
	}
}

// QueryMetricsHandler handles the query_metrics tool call.
func QueryMetricsHandler(ctx context.Context, client *grafana.Client, args map[string]any) (map[string]any, error) {
	// Parse input arguments
	input, err := parseQueryMetricsInput(args)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	slog.Info("Query metrics request",
		"query", input.Query,
		"percentile", input.Percentile,
		"timeRange", input.TimeRange,
		"namespace", input.Namespace,
		"service", input.Service,
		"dataSource", input.DataSource,
	)

	// Get data source if not specified
	dataSource := input.DataSource
	if dataSource == "" {
		defaultDS, err := client.GetDefaultDataSource(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get default data source: %w", err)
		}
		dataSource = defaultDS.UID
		slog.Info("Using default data source", "uid", dataSource, "name", defaultDS.Name)
	}

	// Parse time range
	timeRange, ok := mcpg.ParseTimeRange(input.TimeRange)
	if !ok {
		// Default to 24h
		timeRange, _ = mcpg.ParseTimeRange("24h")
	}

	// Build PromQL query
	promql := BuildPromQL(input)

	slog.Info("Executing PromQL query", "promql", promql, "start", timeRange.Start, "end", timeRange.End, "step", timeRange.Step)

	// Execute query
	result, err := client.DirectQuery(ctx, dataSource, promql, timeRange.Start, timeRange.End, timeRange.Step)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	// Format response
	response := formatQueryResult(result, promql)

	slog.Info("Query completed successfully", "series_count", len(result.Data.Result))

	return response, nil
}

// QueryMetricsInput represents the input for the query_metrics tool.
type QueryMetricsInput struct {
	Query      string
	Percentile float64
	TimeRange  string
	Namespace  string
	Service    string
	DataSource string
	Aggregation string
}

// parseQueryMetricsInput parses the input arguments into QueryMetricsInput.
func parseQueryMetricsInput(args map[string]any) (*QueryMetricsInput, error) {
	input := &QueryMetricsInput{
		TimeRange:   "24h",    // Default
		Aggregation: "avg",    // Default
	}

	if query, ok := args["query"].(string); ok {
		input.Query = query
	}

	if percentile, ok := args["percentile"].(float64); ok {
		input.Percentile = percentile
	}

	if timeRange, ok := args["timeRange"].(string); ok {
		input.TimeRange = timeRange
	}

	if namespace, ok := args["namespace"].(string); ok {
		input.Namespace = namespace
	}

	if service, ok := args["service"].(string); ok {
		input.Service = service
	}

	if dataSource, ok := args["dataSource"].(string); ok {
		input.DataSource = dataSource
	}

	if aggregation, ok := args["aggregation"].(string); ok {
		input.Aggregation = aggregation
	}

	if input.Query == "" {
		return nil, fmt.Errorf("query is required")
	}

	return input, nil
}

// formatQueryResult formats the query response for the MCP tool result.
func formatQueryResult(result *grafana.QueryResponse, promql string) map[string]any {
	series := make([]map[string]any, 0, len(result.Data.Result))

	for _, r := range result.Data.Result {
		values := make([]map[string]any, 0, len(r.Values))
		for _, v := range r.Values {
			if len(v) >= 2 {
				timestamp, _ := v[0].(float64)
				value, _ := v[1].(string)

				var floatVal float64
				fmt.Sscanf(value, "%f", &floatVal)

				values = append(values, map[string]any{
					"timestamp": int64(timestamp),
					"value":     floatVal,
				})
			}
		}

		series = append(series, map[string]any{
			"metric": r.Metric,
			"values": values,
		})
	}

	return map[string]any{
		"promql": promql,
		"status": result.Status,
		"series": series,
		"summary": map[string]any{
			"series_count": len(series),
		},
	}
}

// BuildPromQL builds a PromQL query from the input parameters.
func BuildPromQL(input *QueryMetricsInput) string {
	query := input.Query

	// Add label filters if specified
	if input.Namespace != "" || input.Service != "" {
		// Check if query already has braces (label selectors)
		if len(query) > 0 && query[len(query)-1] == '}' {
			// Query already has label selectors, insert before closing brace
			insertPos := len(query) - 1
			if input.Namespace != "" {
				query = query[:insertPos] + fmt.Sprintf(`,namespace="%s"`, input.Namespace) + query[insertPos:]
			}
			if input.Service != "" {
				query = query[:insertPos] + fmt.Sprintf(`,service="%s"`, input.Service) + query[insertPos:]
			}
		} else if len(query) > 0 {
			// Add label selectors
			labels := []string{}
			if input.Namespace != "" {
				labels = append(labels, fmt.Sprintf(`namespace="%s"`, input.Namespace))
			}
			if input.Service != "" {
				labels = append(labels, fmt.Sprintf(`service="%s"`, input.Service))
			}
			if len(labels) > 0 {
				query = query + "{" + joinStrings(labels, ",") + "}"
			}
		}
	}

	// Apply aggregation
	switch input.Aggregation {
	case "percentile":
		if input.Percentile > 0 {
			// For histogram metrics, use histogram_quantile
			// For duration metrics, use rate with quantile
			query = fmt.Sprintf(`quantile(%.2f, %s)`, input.Percentile/100, query)
		}
	case "avg":
		query = fmt.Sprintf(`avg(%s)`, query)
	case "sum":
		query = fmt.Sprintf(`sum(%s)`, query)
	case "max":
		query = fmt.Sprintf(`max(%s)`, query)
	case "min":
		query = fmt.Sprintf(`min(%s)`, query)
	case "count":
		query = fmt.Sprintf(`count(%s)`, query)
	}

	return query
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
