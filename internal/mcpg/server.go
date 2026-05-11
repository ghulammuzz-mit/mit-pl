package mcpg

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"mit/platform/internal/grafana"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ToolHandlerFunc is the function signature for tool handlers.
type ToolHandlerFunc func(ctx context.Context, client *grafana.Client, args map[string]any) (map[string]any, error)

// Server wraps the MCP server with Grafana client and tool handlers.
type Server struct {
	mcpServer *mcp.Server
	client    *grafana.Client
	handlers  map[string]ToolHandlerFunc
	mu        sync.RWMutex
}

// NewServer creates a new MCP server for Grafana.
func NewServer() (*Server, error) {
	// Create Grafana client
	client := grafana.New()
	if err := client.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate Grafana client: %w", err)
	}

	// Create MCP server
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "mcp-grafana",
		Version: "v0.1.0",
	}, &mcp.ServerOptions{
		// Enable logging
		Logger: slog.Default(),
	})

	srv := &Server{
		mcpServer: mcpServer,
		client:    client,
		handlers:  make(map[string]ToolHandlerFunc),
	}

	slog.Info("MCP Grafana server initialized",
		"name", "mcp-grafana",
		"version", "v0.1.0",
		"grafana_url", client.URL,
	)

	return srv, nil
}

// RegisterTool registers a tool with the server.
func (s *Server) RegisterTool(name string, tool mcp.Tool, handler ToolHandlerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Store the handler
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
		result, err := handler(ctx, s.client, args)
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

// Client returns the Grafana client.
func (s *Server) Client() *grafana.Client {
	return s.client
}

// MCPServer returns the underlying MCP server.
func (s *Server) MCPServer() *mcp.Server {
	return s.mcpServer
}

// Run starts the MCP server with the given transport.
func (s *Server) Run(ctx context.Context, transport mcp.Transport) error {
	slog.Info("Starting MCP Grafana server on transport")
	return s.mcpServer.Run(ctx, transport)
}

// Connect connects the MCP server over the given transport.
func (s *Server) Connect(ctx context.Context, transport mcp.Transport) (*mcp.ServerSession, error) {
	slog.Info("Connecting MCP Grafana server")
	return s.mcpServer.Connect(ctx, transport, nil)
}
