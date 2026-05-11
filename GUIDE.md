# envctl Guide

CLI tool for syncing environment variables between local `.env` files and Infisical.

## Installation

### One-liner (Linux / macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/ghulammuzz-mit/mit-platform/main/scripts/install.sh | sh
```

Installs to `/usr/local/bin/envctl`. Prompts `sudo` if needed.

### One-liner (Windows — PowerShell)

```powershell
irm https://raw.githubusercontent.com/ghulammuzz-mit/mit-platform/main/scripts/install.ps1 | iex
```

Installs to `%LOCALAPPDATA%\Programs\envctl\envctl.exe` and adds to user PATH.

### Verify

```bash
envctl --help
```

### Build from source

```bash
make envctl
# binary output: bin/envctl
```

---

### Platform support

| OS      | Arch  | Binary name                  |
|---------|-------|------------------------------|
| Linux   | amd64 | `envctl-linux-amd64`         |
| Linux   | arm64 | `envctl-linux-arm64`         |
| macOS   | amd64 | `envctl-darwin-amd64`        |
| macOS   | arm64 | `envctl-darwin-arm64`        |
| Windows | amd64 | `envctl-windows-amd64.exe`   |

## Authentication

Credentials are embedded in the binary. No setup required. Connects to `https://us.infisical.com`.

## Global Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--env` | `dev` | Target environment (`dev`, `stg`, `prod`) |
| `--file` | `.env` | Local env file path |
| `--yes` | `false` | Skip confirmation prompts |

---

## Commands

### `list`

List all projects and folders in Infisical.

```bash
envctl list
envctl list --env prod
```

**Output:**
```
=== Projects ===
- my-app  abc123
- another-app  def456

=== Folders ===
- backend
- frontend
```

> Note: Only lists folders from the first project returned.

---

### `push`

Upload local `.env` to Infisical.

```bash
envctl push
envctl push --env prod
envctl push --env stg --file .env.staging
envctl push --env dev --yes   # skip confirmation
```

**Flow:**
1. Prompts to select a project (interactive)
2. Prompts to select a folder/app (interactive)
3. Reads local `--file`
4. Asks confirmation unless `--yes` passed
5. Deletes all existing secrets in that path
6. Uploads all key/value pairs from local file

> **Warning:** Push is destructive — all existing secrets in the selected path are deleted before new ones are created. Ensure your local file is complete before pushing.

---

### `pull`

Download secrets from Infisical and write to local file.

```bash
envctl pull
envctl pull --env prod
envctl pull --env stg --file .env.staging
```

**Flow:**
1. Prompts to select a project (interactive)
2. Prompts to select a folder/app (interactive)
3. Fetches all secrets from Infisical
4. Overwrites local `--file` with fetched secrets

> **Warning:** Pull overwrites the local file completely. Back up any local changes before pulling.

---

## Common Workflows

### First-time setup: pull dev secrets

```bash
envctl pull --env dev
```

### Deploy: push staging secrets

```bash
envctl push --env stg --file .env.staging --yes
```

### Sync prod to local (read-only inspect)

```bash
envctl pull --env prod --file .env.prod.local
```

### Interactive session (default)

Running any command without `--yes` prompts for project and folder selection:

```
=== Select Project ===
[1] my-app
[2] another-app
Choose: 1

=== Select App (Folder) ===
[1] backend
[2] frontend
Choose: 1
```

---

## Notes

- `push` deletes then re-creates all secrets — partial uploads not supported
- `pull` format: `KEY=VALUE` per line, no comments preserved
- `--env` must match an environment that exists in Infisical (e.g., `dev`, `stg`, `prod`)
- `list` only shows folders from the first project — use `push`/`pull` for other projects
