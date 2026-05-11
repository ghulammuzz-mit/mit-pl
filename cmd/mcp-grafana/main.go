package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"mit/platform/internal/mcpg"
	"mit/platform/internal/mcpg/tools"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		slog.Warn("Failed to load .env file", "error", err)
	}

	// Set up structured logging
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	slog.Info("Starting MCP Grafana server")

	// Validate environment variables
	if os.Getenv("GRAFANA_URL") == "" {
		slog.Error("GRAFANA_URL environment variable is required")
		os.Exit(1)
	}
	if os.Getenv("GRAFANA_API_KEY") == "" {
		slog.Error("GRAFANA_API_KEY environment variable is required")
		os.Exit(1)
	}

	// Create MCP server
	server, err := mcpg.NewServer()
	if err != nil {
		slog.Error("Failed to create MCP server", "error", err)
		os.Exit(1)
	}

	// Register tools
	server.RegisterTool("query_metrics", tools.QueryMetricsTool(), tools.QueryMetricsHandler)
	server.RegisterTool("list_dashboards", tools.ListDashboardsTool(), tools.ListDashboardsHandler)
	server.RegisterTool("list_data_sources", tools.ListDataSourcesTool(), tools.ListDataSourcesHandler)

	// Set up context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		slog.Info("Received shutdown signal")
		cancel()
	}()

	// Create stdio transport
	transport := &mcp.StdioTransport{}

	// Start server
	slog.Info("MCP Grafana server running on stdio transport")
	if err := server.Run(ctx, transport); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}

	slog.Info("MCP Grafana server stopped")
}
