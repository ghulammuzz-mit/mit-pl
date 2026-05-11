#!/bin/bash

# Test script for MCP Infisical server using proper JSON-RPC

export INFISICAL_UNIVERSAL_AUTH_CLIENT_ID="${INFISICAL_UNIVERSAL_AUTH_CLIENT_ID:-}"
export INFISICAL_UNIVERSAL_AUTH_CLIENT_SECRET="${INFISICAL_UNIVERSAL_AUTH_CLIENT_SECRET:-}"
export INFISICAL_HOST_URL="${INFISICAL_HOST_URL:-https://app.infisical.com}"

echo "=== MCP Infisical Server Test ==="
echo ""

if [ -z "$INFISICAL_UNIVERSAL_AUTH_CLIENT_ID" ] || [ -z "$INFISICAL_UNIVERSAL_AUTH_CLIENT_SECRET" ]; then
	echo "⚠️  Warning: INFISICAL_UNIVERSAL_AUTH_CLIENT_ID and/or INFISICAL_UNIVERSAL_AUTH_CLIENT_SECRET not set"
	echo "   Set these environment variables to test the server"
	echo ""
fi

# Create a temporary file with the JSON-RPC requests
cat > /tmp/mcp_infisical_test.json << 'EOF'
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}
{"jsonrpc":"2.0","id":2,"method":"notifications/initialized"}
{"jsonrpc":"2.0","id":3,"method":"tools/list"}
EOF

echo "Testing MCP server with JSON-RPC requests..."
echo ""
go run ./cmd/mcp-infisical < /tmp/mcp_infisical_test.json 2>&1 | head -50

echo ""
echo "=== Test Complete ==="
