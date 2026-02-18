# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Go CLI tool that migrates content from Obsidian vaults to Notion workspaces. It reads Obsidian Markdown files (including frontmatter, wikilinks, and attachments) and creates/updates corresponding Notion pages via the Notion API.

## Common Commands

```bash
# Initialize module (first time only)
go mod init github.com/kfi/obsidian-to-notion

# Build
go build ./...
go build -o obsidian-to-notion ./cmd/obsidian-to-notion

# Run
./obsidian-to-notion --vault /path/to/vault --token <notion-token> --database <db-id>

# Test
go test ./...
go test ./... -v
go test ./internal/parser/...        # single package
go test -run TestFunctionName ./...  # single test

# Lint (requires golangci-lint)
golangci-lint run

# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Tidy dependencies
go mod tidy
```

## Architecture

```
cmd/obsidian-to-notion/   # main package and CLI entry point (cobra or flag-based)
internal/
  obsidian/               # Vault reading: file walking, Markdown parsing, frontmatter, wikilinks
  notion/                 # Notion API client wrapper (pages, blocks, databases)
  converter/              # Maps Obsidian AST/structures to Notion block types
  sync/                   # Orchestrates the migration: tracks state, handles updates vs creates
  config/                 # Config loading (env vars, config file, CLI flags)
```

### Data Flow

1. **Obsidian layer** walks the vault directory and parses each `.md` file — extracting YAML frontmatter, body content, wikilinks (`[[note]]`), and tags.
2. **Converter** transforms Obsidian Markdown into Notion block objects (paragraphs, headings, bullets, callouts, etc.).
3. **Notion layer** wraps the official Notion SDK and handles rate limiting, pagination, and retries.
4. **Sync layer** determines whether to create or update a page (keyed on a stored Notion page ID or a title match) and applies changes incrementally.

### Key Dependencies

- Notion API: stdlib `net/http` — no SDK, calls the REST API directly (`baseURL = https://api.notion.com/v1`, `Notion-Version: 2022-06-28`)
- Markdown parsing: `github.com/yuin/goldmark`
- CLI framework: `github.com/spf13/cobra`
- YAML frontmatter: `github.com/adrg/frontmatter` or `gopkg.in/yaml.v3`

## Configuration

The tool is configured via a combination of CLI flags and environment variables:

| Flag | Env var | Description |
|------|---------|-------------|
| `--vault` | `OBSIDIAN_VAULT` | Path to the Obsidian vault root |
| `--token` | `NOTION_TOKEN` | Notion integration token |
| `--database` | `NOTION_DATABASE_ID` | Target Notion database ID |
| `--config` | — | Path to a YAML config file |

## Notion API Notes

- Notion blocks have a maximum nesting depth of 3; deeper Obsidian nesting must be flattened.
- Rich text content is limited to 2000 characters per block; long paragraphs must be split.
- API requests are rate-limited to ~3 req/s; the Notion client must implement exponential backoff.
- Page content cannot be partially updated — replacing blocks requires deleting existing children and re-appending.
