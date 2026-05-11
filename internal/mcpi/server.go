package mcpi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"

	infisical "github.com/infisical/go-sdk"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server wraps MCP server with Infisical client and tool handlers.
type Server struct {
	mcpServer *mcp.Server
	client    infisical.InfisicalClientInterface
	handlers  map[string]ToolHandlerFunc
	mu        sync.RWMutex
}

// NewServer creates a new MCP server for Infisical.
func NewServer() (*Server, error) {
	// Get environment variables for authentication
	clientID := os.Getenv("INFISICAL_UNIVERSAL_AUTH_CLIENT_ID")
	clientSecret := os.Getenv("INFISICAL_UNIVERSAL_AUTH_CLIENT_SECRET")
	siteURL := os.Getenv("INFISICAL_HOST_URL")
	if siteURL == "" {
		siteURL = "https://us.infisical.com"
	}

	if clientID == "" {
		return nil, &ConfigError{Field: "INFISICAL_UNIVERSAL_AUTH_CLIENT_ID", Message: "Infisical Universal Auth Client ID is required"}
	}
	if clientSecret == "" {
		return nil, &ConfigError{Field: "INFISICAL_UNIVERSAL_AUTH_CLIENT_SECRET", Message: "Infisical Universal Auth Client Secret is required"}
	}

	// Create Infisical client
	ctx := context.Background()
	client := infisical.NewInfisicalClient(ctx, infisical.Config{
		SiteUrl:          siteURL,
		AutoTokenRefresh: true,
	})

	// Authenticate with Universal Auth
	_, err := client.Auth().UniversalAuthLogin(clientID, clientSecret)
	if err != nil {
		return nil, fmt.Errorf("Infisical authentication failed: %w", err)
	}

	// Create MCP server
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "mcp-infisical",
		Version: "v0.1.0",
	}, &mcp.ServerOptions{
		Logger: slog.Default(),
	})

	srv := &Server{
		mcpServer: mcpServer,
		client:    client,
		handlers:  make(map[string]ToolHandlerFunc),
	}

	slog.Info("MCP Infisical server initialized",
		"name", "mcp-infisical",
		"version", "v0.1.0",
		"site_url", siteURL,
	)

	return srv, nil
}

// RegisterTool registers a tool with the server.
func (s *Server) RegisterTool(name string, tool mcp.Tool, handler ToolHandlerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Store handler
	s.handlers[name] = handler

	// Wrap handler to match mcp.ToolHandler signature
	wrappedHandler := mcp.ToolHandler(func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Unmarshal arguments
		var args map[string]any
		if len(req.Params.Arguments) > 0 {
			if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
				return &mcp.CallToolResult{
					IsError: true,
					Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("failed to unmarshal arguments: %v", err)}},
				}, nil
			}
		}

		slog.Info("Handling tool call", "tool", name, "args", args)
		result, err := handler(args)
		if err != nil {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
			}, nil
		}

		// Convert result to JSON content
		resultJSON, _ := json.Marshal(result)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(resultJSON)}},
		}, nil
	})

	// Add tool to MCP server
	s.mcpServer.AddTool(&tool, wrappedHandler)

	slog.Info("Registered tool", "name", name)
}

// Client returns the Infisical client.
func (s *Server) Client() infisical.InfisicalClientInterface {
	return s.client
}

// MCPServer returns the underlying MCP server.
func (s *Server) MCPServer() *mcp.Server {
	return s.mcpServer
}

// Run starts the MCP server with the given transport.
func (s *Server) Run(ctx context.Context, transport mcp.Transport) error {
	slog.Info("Starting MCP Infisical server on transport")
	return s.mcpServer.Run(ctx, transport)
}

// ConfigError represents a configuration error.
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return e.Field + ": " + e.Message
}
