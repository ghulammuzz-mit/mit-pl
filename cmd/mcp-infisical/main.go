package main

import (
	"context"
	"log/slog"
	"mit/platform/internal/mcpi"
	"mit/platform/internal/mcpi/tools"

	"github.com/joho/godotenv"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	// Create MCP server
	server, err := mcpi.NewServer()
	if err != nil {
		slog.Error("Failed to create MCP server", "error", err)
		return
	}

	// Register all tools
	tools.RegisterCreateSecret(server)
	tools.RegisterDeleteSecret(server)
	tools.RegisterUpdateSecret(server)
	tools.RegisterListSecrets(server)
	tools.RegisterGetSecret(server)
	tools.RegisterCreateProject(server)
	tools.RegisterCreateEnvironment(server)
	tools.RegisterCreateFolder(server)

	slog.Info("MCP Infisical server starting",
		"tools_count", 8,
	)

	// Create stdio transport and run
	transport := &mcp.StdioTransport{}
	ctx := context.Background()

	if err := server.Run(ctx, transport); err != nil {
		slog.Error("Server error", "error", err)
	}
}
