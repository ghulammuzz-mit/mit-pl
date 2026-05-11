# mit-platform

Go monorepo containing CLI tools and services for infrastructure management.

## Tools

| Tool | Description |
|------|-------------|
| [envctl](#envctl) | Sync `.env` files with Infisical secrets manager |
| [mcp-grafana](#mcp-grafana) | MCP server for querying Prometheus metrics via Grafana |
| [mcp-infisical](#mcp-infisical) | MCP server for managing Infisical secrets |
| [dns](#dns) | HTTP service for managing Cloudflare DNS records |
| [hpa-metrics](#hpa-metrics) | Kubernetes HPA planning utility |

---

## envctl

CLI tool for syncing environment variables between local `.env` files and Infisical.

### Install

**Linux / macOS**
```bash
curl -fsSL https://raw.githubusercontent.com/ghulammuzz-mit/mit-pl/main/scripts/install.sh | sh
```

**Windows (PowerShell)**
```powershell
irm https://raw.githubusercontent.com/ghulammuzz-mit/mit-pl/main/scripts/install.ps1 | iex
```

**Build from source**
```bash
make envctl
# output: bin/envctl
```

**Platform support**

| OS      | Arch  | Binary                       |
|---------|-------|------------------------------|
| Linux   | amd64 | `envctl-linux-amd64`         |
| Linux   | arm64 | `envctl-linux-arm64`         |
| macOS   | amd64 | `envctl-darwin-amd64`        |
| macOS   | arm64 | `envctl-darwin-arm64`        |
| Windows | amd64 | `envctl-windows-amd64.exe`   |

### Configuration

```bash
INFISICAL_UNIVERSAL_AUTH_CLIENT_ID=your-client-id
INFISICAL_UNIVERSAL_AUTH_CLIENT_SECRET=your-client-secret
```

Load from `.env` file or set in shell. Copy `.env.example` to get started:
```bash
cp .env.example .env
```

### Commands

```bash
envctl list                                      # list projects and folders
envctl pull                                      # pull secrets to .env
envctl pull --env prod                           # pull from prod environment
envctl pull --env stg --file .env.staging        # pull to custom file
envctl push                                      # push .env to Infisical
envctl push --env prod --yes                     # push to prod, skip confirmation
envctl push --env stg --file .env.staging        # push from custom file
```

**Global flags**

| Flag | Default | Description |
|------|---------|-------------|
| `--env` | `dev` | Target environment (`dev`, `stg`, `prod`) |
| `--file` | `.env` | Local env file path |
| `--yes` | `false` | Skip confirmation prompt |

> **Warning:** `push` deletes all existing secrets in the target path before uploading.
> `pull` overwrites the local file completely. Back up first.

### Interactive flow

All commands prompt for project and folder selection:

```
=== Select Project ===
[1] my-app
[2] backend-api
Choose: 1

=== Select App (Folder) ===
[1] backend
[2] frontend
Choose: 1
```

---

## mcp-grafana

MCP (Model Context Protocol) server for querying Prometheus metrics and Grafana dashboards via Claude Code.

### Configuration

```bash
GRAFANA_URL=https://your-grafana.example.com
GRAFANA_API_KEY=your-grafana-api-key
```

### Claude Code setup

```json
{
  "mcpServers": {
    "grafana": {
      "command": "/path/to/bin/mcp-grafana",
      "env": {
        "GRAFANA_URL": "https://your-grafana.example.com",
        "GRAFANA_API_KEY": "your-grafana-api-key"
      }
    }
  }
}
```

### Tools

#### `query_metrics`
Query Prometheus metrics with filters and aggregations.

| Parameter | Required | Description |
|-----------|----------|-------------|
| `query` | yes | PromQL query or metric name |
| `timeRange` | no | `1h`, `6h`, `24h`, `7d`, `30d` (default: `24h`) |
| `namespace` | no | Kubernetes namespace filter |
| `service` | no | Service name filter |
| `aggregation` | no | `avg`, `sum`, `max`, `min`, `count`, `percentile` |
| `percentile` | no | Percentile value (e.g., `95` for p95) |
| `dataSource` | no | Grafana data source UID |

#### `list_dashboards`
List Grafana dashboards with optional search and tag filtering.

#### `list_data_sources`
List Grafana data sources, optionally filtered by type (e.g., `prometheus`).

### Build & run

```bash
make mcp-grafana
# output: bin/mcp-grafana
```

---

## mcp-infisical

MCP server for managing Infisical secrets via Claude Code.

### Configuration

```bash
INFISICAL_UNIVERSAL_AUTH_CLIENT_ID=your-client-id
INFISICAL_UNIVERSAL_AUTH_CLIENT_SECRET=your-client-secret
INFISICAL_HOST_URL=https://us.infisical.com   # optional
```

### Claude Code setup

```json
{
  "mcpServers": {
    "infisical": {
      "command": "/path/to/bin/mcp-infisical",
      "env": {
        "INFISICAL_UNIVERSAL_AUTH_CLIENT_ID": "your-client-id",
        "INFISICAL_UNIVERSAL_AUTH_CLIENT_SECRET": "your-client-secret"
      }
    }
  }
}
```

### Tools

| Tool | Description |
|------|-------------|
| `create_secret` | Create a new secret |
| `get_secret` | Get a single secret by key |
| `update_secret` | Update an existing secret |
| `delete_secret` | Delete a secret |
| `list_secrets` | List all secrets in a path |
| `create_folder` | Create a folder |

### Build & run

```bash
make mcp-infisical
# output: bin/mcp-infisical
```

---

## dns

HTTP service for managing Cloudflare DNS records.

### Configuration

```bash
PORT=8080
CF_API_TOKEN=your-cloudflare-api-token
CF_ZONE_ID=your-cloudflare-zone-id
```

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/dns` | List all DNS records |
| `POST` | `/dns` | Create a DNS record |
| `PUT` | `/dns/{id}` | Update a DNS record by ID |
| `DELETE` | `/dns/{name}` | Delete a DNS record by name |

### Build & run

```bash
make dns
bin/dns
```

---

## hpa-metrics

Analyzes Prometheus metrics and generates Kubernetes HPA recommendations.

### Configuration

```bash
GRAFANA_URL=https://your-grafana.example.com
GRAFANA_API_KEY=your-grafana-api-key
```

### Usage

```bash
make run-hpa-metrics
```

**Output:**
- `hpa-planning/report.md` — recommendations in markdown
- `hpa-planning/metrics.json` — raw metrics data

---

## Development

### Prerequisites

- Go 1.24+
- `make`

### Build all

```bash
make build
# outputs: bin/envctl bin/dns bin/mcp-grafana bin/mcp-infisical bin/hpa-metrics
```

### Environment setup

```bash
cp .env.example .env
# fill in .env with your credentials
```

### Cross-compile envctl for all platforms

```bash
make release-envctl
# outputs: dist/envctl-{os}-{arch}
```
