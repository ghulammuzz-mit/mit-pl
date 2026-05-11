#!/bin/bash

# Test script for MCP Grafana server using proper JSON-RPC

export GRAFANA_URL="${GRAFANA_URL:-}"
export GRAFANA_API_KEY="${GRAFANA_API_KEY:-}"

if [ -z "$GRAFANA_URL" ] || [ -z "$GRAFANA_API_KEY" ]; then
	echo "Error: GRAFANA_URL and GRAFANA_API_KEY must be set"
	exit 1
fi

echo "=== MCP Grafana Server Test ==="
echo ""

# Create a temporary file with the JSON-RPC requests
cat > /tmp/mcp_test_input.json << 'EOF'
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}
{"jsonrpc":"2.0","id":2,"method":"notifications/initialized"}
{"jsonrpc":"2.0","id":3,"method":"tools/list"}
EOF

echo "Testing MCP server with JSON-RPC requests..."
echo ""
go run ./cmd/mcp-grafana < /tmp/mcp_test_input.json 2>&1 | head -50

echo ""
echo "=== Test Complete ==="
