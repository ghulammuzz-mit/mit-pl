package mcpi

// ToolHandlerFunc is the function signature for Infisical MCP tool handlers.
// It returns a map[string]any result and an error.
type ToolHandlerFunc func(args map[string]any) (map[string]any, error)
